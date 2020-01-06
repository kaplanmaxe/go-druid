package dsql

import "database/sql/driver"

type result struct{}

// LastInsertId is a noop
func (r *result) LastInsertId() (id int64, err error) {
	return id, driver.ErrSkip
}

// RowsAffected is a noop as this driver does not support inserts
func (r *result) RowsAffected() (rows int64, err error) {
	return rows, driver.ErrSkip
}
