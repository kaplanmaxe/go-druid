package druid

type result struct{}

func (r *result) LastInsertId() (id int64, err error) {
	return
}

func (r *result) RowsAffected() (rows int64, err error) {
	return
}
