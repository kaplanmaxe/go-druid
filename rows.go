package druid

import (
	"database/sql/driver"
	"errors"
	"io"
	"log"
)

type resultSet struct {
	rows        [][]field
	columnNames []string
	currentRow  int
}

type rows struct {
	conn      *connection
	resultSet resultSet
}

func (r *rows) Columns() (cols []string) {
	cols = r.resultSet.columnNames
	return
}

func (r *rows) Close() (err error) {
	return
}

func (r *rows) Next(dest []driver.Value) (err error) {
	if r.HasNextResultSet() == false {
		return io.EOF
	}

	data := r.resultSet.rows[r.resultSet.currentRow]
	if len(data) != len(dest) {
		return errors.New("druid: number of refs passed to scan does not match column count")
	}
	for i := range dest {
		switch data[i].Type.Name() {
		// TODO: add time.Time and []byte
		case "bool":
			dest[i] = data[i].Value.Interface().(bool)
		case "string":
			dest[i] = data[i].Value.Interface().(string)
		case "int":
			dest[i] = data[i].Value.Interface().(int)
		case "int64":
			dest[i] = data[i].Value.Interface().(int64)
		case "float64":
			dest[i] = data[i].Value.Interface().(float64)
		default:
			log.Fatal("druid: can't scan type ", data[i].Type.Name())
		}
	}
	r.NextResultSet()

	return
}

// HasNextResultSet implements driver.RowsNextResultSet
func (r *rows) HasNextResultSet() bool {
	if r.resultSet.currentRow == len(r.resultSet.rows) {
		return false
	}
	return true
}

// NextResultSet implements driver.RowsNextResultSet
func (r *rows) NextResultSet() error {
	r.resultSet.currentRow++
	if r.HasNextResultSet() == false {
		return io.EOF
	}
	return nil
}
