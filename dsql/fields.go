package dsql

import (
	"reflect"
)

type field struct {
	Value reflect.Value
	Type  reflect.Type
}
