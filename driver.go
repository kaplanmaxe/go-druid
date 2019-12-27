package druid

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

// Driver is a struct meant to be returned and used with database/sql
type Driver struct{}

func init() {
	sql.Register("druid", &Driver{})
}

// Open opens a new connection
func (d *Driver) Open(dsn string) (driver.Conn, error) {
	cfg := ParseDSN(dsn)
	conn := &connector{
		Cfg: cfg,
	}
	return conn.Connect(context.Background())
}
