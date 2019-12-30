package druid

import (
	"database/sql/driver"
	"testing"
)

func TestLastInsertId(t *testing.T) {
	stmt := &statementNoop{}
	res, _ := stmt.Exec([]driver.Value{})
	_, err := res.LastInsertId()
	if err != driver.ErrSkip {
		t.Fatal("Expected LastInsertId to be unimplemented but it is")
	}
}

func TestRowsAffected(t *testing.T) {
	stmt := &statementNoop{}
	res, _ := stmt.Exec([]driver.Value{})
	_, err := res.RowsAffected()
	if err != driver.ErrSkip {
		t.Fatal("Expected RowsAffected to be unimplemented but it is")
	}
}
