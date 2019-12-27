package druid_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kaplanmaxe/druid"
)

var cfg druid.Config = druid.Config{
	BrokerAddr:   "localhost:8082",
	PingEndpoint: "/status/health",
	// User:         "druidUsername",
	// Passwd:       "druidPassword",
}

func startMockServer(handler http.HandlerFunc) (ts *httptest.Server) {
	ts = httptest.NewServer(handler)
	return
}

func TestConnect(t *testing.T) {
	_, err := sql.Open("druid", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPing(t *testing.T) {
	ts := startMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()
	cfg.BrokerAddr = strings.Replace(ts.URL, "http://", "", 1)
	db, err := sql.Open("druid", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
