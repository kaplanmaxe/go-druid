package druid

import (
	"context"
	"database/sql/driver"
	"net/http"
)

type connector struct {
	Cfg *Config
}

// Connect implements db.Connector and sets up an http client to druid's sql endpoint
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	client := &http.Client{}
	connection := &connection{
		Client:    client,
		Cfg:       c.Cfg,
		closeCh:   make(chan struct{}, 1),
		watcherCh: make(chan context.Context, 1),
		errorCh:   make(chan error),
		resultsCh: make(chan []byte),
		requestCh: make(chan *http.Request),
	}
	connection.startRequestPipeline()
	return connection, nil
}

// Driver implements db.Connector and returns a druid driver
func (c *connector) Driver() (d driver.Driver) {
	return &Driver{}
}
