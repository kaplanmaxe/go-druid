package druid

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrPinging is an error returned when health check endpoint returns a non 200.
	ErrPinging = errors.New("druid: error fetching health info from druid")
)

type connection struct {
	Client *http.Client
	Cfg    *Config
}

// Prepare implements db.Conn.Prepare and returns a noop statement
func (c *connection) Prepare(stmt string) (driver.Stmt, error) {
	return &statementNoop{}, nil
}

// Close closes a connection.
// TODO: implement
func (c *connection) Close() (err error) {
	return
}

// Begin implements db.Conn.Prepare and is a noop
func (c *connection) Begin() (tx driver.Tx, err error) {
	tx = &transactionNoop{}
	return
}

// Ping implmements db.conn.Prepare and hits the health endpoint of a broker
func (c *connection) Ping(ctx context.Context) (err error) {
	res, err := c.Client.Get(fmt.Sprintf("%s%s", c.Cfg.BrokerAddr, c.Cfg.PingEndpoint))
	if err != nil {
		err = ErrPinging
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("druid: got %d status code from %s", res.StatusCode, c.Cfg.PingEndpoint)
	}
	return
}
