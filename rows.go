package druid

import (
	"database/sql/driver"
	"io"
	"log"
)

type resultSet struct {
	columns     []field
	columnNames []string
	currentCol  int
	done        bool
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
	if r.resultSet.currentCol == len(r.resultSet.columns) {
		return io.EOF
	}

	// log.Fatal(data)
	for i := range dest {
		data := r.resultSet.columns[r.resultSet.currentCol]
		switch data.Type.Name() {
		case "string":
			dest[i] = data.Value.Interface().(string)
		}
		r.resultSet.currentCol++
	}

	return
}

func (r *rows) HasNextResultSet(b bool) {
	log.Fatal("called")
	return
}
