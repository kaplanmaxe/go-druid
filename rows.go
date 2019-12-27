package druid

import "database/sql/driver"

type rows struct{}

func (r *rows) Columns() (cols []string) {
	return
}

func (r *rows) Close() (err error) {
	return
}

func (r *rows) Next(dest []driver.Value) (err error) {
	return
}
