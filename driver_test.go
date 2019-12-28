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

var mockQueryResults = `[["__time","channel"],["2015-09-12T00:46:58.771Z","#en.wikipedia"],["2015-09-12T00:47:00.496Z","#ca.wikipedia"],["2015-09-12T00:47:05.474Z","#en.wikipedia"],["2015-09-12T00:47:08.770Z","#vi.wikipedia"],["2015-09-12T00:47:11.862Z","#vi.wikipedia"],["2015-09-12T00:47:13.987Z","#vi.wikipedia"],["2015-09-12T00:47:17.009Z","#ca.wikipedia"],["2015-09-12T00:47:19.591Z","#en.wikipedia"],["2015-09-12T00:47:21.578Z","#en.wikipedia"],["2015-09-12T00:47:25.821Z","#vi.wikipedia"]]`

func startMockServer(handler http.HandlerFunc) (ts *httptest.Server, url string) {
	ts = httptest.NewServer(handler)
	url = strings.Replace(ts.URL, "http://", "", 1)
	return
}

func TestConnect(t *testing.T) {
	_, err := sql.Open("druid", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPing(t *testing.T) {
	ts, url := startMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer ts.Close()
	cfg.BrokerAddr = url
	db, err := sql.Open("druid", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	ts, url := startMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockQueryResults))
		w.Header().Add("Content-Type", "application/json")
	})
	defer ts.Close()
	cfg.BrokerAddr = url
	db, err := sql.Open("druid", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := db.Query("SELECT __time, channel FROM \"example\" LIMIT 10")
	if err != nil {
		t.Fatal(err)
	}
	var channel string
	var time string
	var channels []string
	var times []string
	for rows.Next() {
		err := rows.Scan(&time, &channel)
		if err != nil {
			t.Error(err)
		}
		channels = append(channels, channel)
		times = append(times, time)
	}
	if len(times) == 0 || len(channels) == 0 {
		t.Error("Did not fetch results properly")
	}
}
