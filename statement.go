package druid

import (
	"context"
	"database/sql/driver"
)

type statementNoop struct{}

// Close is a noop as druid is not an OLTP DB
func (s *statementNoop) Close() (err error) {
	return driver.ErrSkip
}

// NumInput is a noop as druid is not an OLTP DB
func (s *statementNoop) NumInput() (num int) {
	return
}

// ExecContext is a noop as druid is not an OLTP DB
func (s *statementNoop) ExecContext(ctx context.Context, args []driver.Value) (driver.Result, error) {
	return &result{}, driver.ErrSkip
}

// Exec is a noop as druid is not an OLTP DB
func (s *statementNoop) Exec(args []driver.Value) (driver.Result, error) {
	return &result{}, driver.ErrSkip
}

// Query is a noop as druid is not an OLTP DB
func (s *statementNoop) Query(args []driver.Value) (driver.Rows, error) {
	return &rows{}, driver.ErrSkip
}
