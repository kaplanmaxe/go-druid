package druid

import (
	"database/sql/driver"
	"errors"
)

func namedValuesToValues(namedValues []driver.NamedValue) (values []driver.Value, err error) {
	for _, val := range namedValues {
		if len(val.Name) > 0 {
			return values, errors.New("druid: named values not supported")
		}
		values = append(values, val)
	}
	return
}
