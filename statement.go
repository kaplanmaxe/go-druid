package druid

import (
	"context"
	"database/sql/driver"
)

type statementNoop struct{}

func (s *statementNoop) Close() (err error) {
	return
}

func (s *statementNoop) NumInput() (num int) {
	return
}

func (s *statementNoop) ExecContext(ctx context.Context, args []driver.Value) (driver.Result, error) {
	return &result{}, nil
}

func (s *statementNoop) Exec(args []driver.Value) (driver.Result, error) {
	return &result{}, nil
}

func (s *statementNoop) Query(args []driver.Value) (driver.Rows, error) {
	return &rows{}, nil
}
