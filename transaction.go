package druid

type transactionNoop struct{}

func (t *transactionNoop) Commit() (err error) {
	return
}

func (t *transactionNoop) Rollback() (err error) {
	return
}
