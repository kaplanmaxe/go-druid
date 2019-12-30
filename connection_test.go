package druid_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/kaplanmaxe/druid"
)

func TestPrepare(t *testing.T) {
	cfg := druid.Config{
		BrokerAddr:   "localhost:8082",
		PingEndpoint: "/status/health",
	}
	db, err := sql.Open("druid", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Prepare("SELECT * FROM example")
	if err != driver.ErrSkip {
		t.Fatal("Expected prepare to be unimplemented but it is")
	}
}

func TestBegin(t *testing.T) {
	cfg := druid.Config{
		BrokerAddr:   "localhost:8082",
		PingEndpoint: "/status/health",
	}
	db, err := sql.Open("druid", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Begin()
	if err != driver.ErrSkip {
		t.Fatal("Expected begin to be unimplemented but it is")
	}
}
