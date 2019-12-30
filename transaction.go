package druid

import "database/sql/driver"

type transactionNoop struct{}

// Commit is a noop that as druid is not an OLTP DB
func (t *transactionNoop) Commit() (err error) {
	return driver.ErrSkip
}

// Rollback is a noob as druid is not an OLTP DB
func (t *transactionNoop) Rollback() (err error) {
	return driver.ErrSkip
}
