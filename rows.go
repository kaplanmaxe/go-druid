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

	for i := range dest {
		data := r.resultSet.columns[r.resultSet.currentCol]
		switch data.Type.Name() {
		case "bool":
			dest[i] = data.Value.Interface().(bool)
		case "string":
			dest[i] = data.Value.Interface().(string)
		case "int":
			dest[i] = data.Value.Interface().(int)
		case "int8":
			dest[i] = data.Value.Interface().(int8)
		case "int16":
			dest[i] = data.Value.Interface().(int16)
		case "int32":
			dest[i] = data.Value.Interface().(int32)
		case "int64":
			dest[i] = data.Value.Interface().(int64)
		case "uint":
			dest[i] = data.Value.Interface().(uint)
		case "uint8":
			dest[i] = data.Value.Interface().(uint8)
		case "uint16":
			dest[i] = data.Value.Interface().(uint16)
		case "uint32":
			dest[i] = data.Value.Interface().(uint32)
		case "uint64":
			dest[i] = data.Value.Interface().(uint64)
		case "uintptr":
			// Don't think it's needed but can't hurt
			dest[i] = data.Value.Interface().(uintptr)
		case "rune":
			dest[i] = data.Value.Interface().(rune)
		case "float32":
			dest[i] = data.Value.Interface().(float32)
		case "float64":
			dest[i] = data.Value.Interface().(float64)
		case "complex64":
			dest[i] = data.Value.Interface().(complex64)
		case "complex128":
			dest[i] = data.Value.Interface().(complex128)
		default:
			log.Fatal("druid: can't scan type ", data.Type.Name())
		}
		r.resultSet.currentCol++
	}

	return
}

func (r *rows) HasNextResultSet(b bool) {
	log.Fatal("called")
	return
}
