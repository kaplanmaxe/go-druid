package druid_test

import (
	"database/sql"
	"encoding/json"
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

// type mockResults struct {
// 	BoolTest    bool    `json:"bool_test"`
// 	StringTest  string  `json:"string_test"`
// 	IntTest     int     `json:"int_test"`
// 	Int8Test    int8    `json:"int8_test"`
// 	Int16Test   int16   `json:"int16_test"`
// 	Int32Test   int32   `json:"int32_test"`
// 	Int64Test   int64   `json:"int64_test"`
// 	UintTest    uint    `json:"uint_test"`
// 	Uint8Test   uint8   `json:"uint8_test"`
// 	Uint16Test  uint16  `json:"uint16_test"`
// 	Uint32Test  uint32  `json:"uint32_test"`
// 	Uint64Test  uint64  `json:"uint64_test"`
// 	Float32Test float32 `json:"float32_test"`
// 	Float64Test float64 `json:"float64_test"`
// }

// TODO: better method for constructing tests
var mockQueryResults = `[["__time","added","channel"],["2015-09-12T00:46:58.771Z",36,"#en.wikipedia"],["2015-09-12T00:47:00.496Z",17,"#ca.wikipedia"],["2015-09-12T00:47:05.474Z",0,"#en.wikipedia"],["2015-09-12T00:47:08.770Z",18,"#vi.wikipedia"],["2015-09-12T00:47:11.862Z",18,"#vi.wikipedia"],["2015-09-12T00:47:13.987Z",18,"#vi.wikipedia"],["2015-09-12T00:47:17.009Z",0,"#ca.wikipedia"],["2015-09-12T00:47:19.591Z",345,"#en.wikipedia"],["2015-09-12T00:47:21.578Z",121,"#en.wikipedia"],["2015-09-12T00:47:25.821Z",18,"#vi.wikipedia"]]`

func startMockServer(handler http.HandlerFunc) (ts *httptest.Server, url string) {
	ts = httptest.NewServer(handler)
	url = strings.Replace(ts.URL, "http://", "", 1)
	return
}

func constructMockResults(header []interface{}, rows [][]interface{}) (b []byte, err error) {
	var mockResults [][]interface{}
	mockResults = append(mockResults, header)
	for _, val := range rows {
		mockResults = append(mockResults, val)
	}
	b, err = json.Marshal(mockResults)
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
	header := []interface{}{"__time", "added", "channel"}
	mockRows := [][]interface{}{{"2015-09-12T00:46:58.771Z", 36, "#en.wikipedia"}, {"2015-09-12T00:46:58.772Z", 76, "#ca.wikipedia"}}
	output, _ := constructMockResults(header, mockRows)
	ts, url := startMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.Write(output)
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
	rows, err := db.Query("SELECT __time, added, channel FROM \"wikiticker-2015-09-12-sampled\" LIMIT 10")
	if err != nil {
		t.Fatal(err)
	}
	var channel string
	var time string
	var added int
	var channels []string
	var times []string
	var addeds []int
	row := 0
	for rows.Next() {
		err := rows.Scan(&time, &added, &channel)
		if err != nil {
			t.Error(err)
		}
		if time != mockRows[row][0] {
			t.Fatalf("Expecting %v got %v", mockRows[row][0], time)
		}
		if added != mockRows[row][1] {
			t.Fatalf("Expecting %v got %v", mockRows[row][1], added)
		}
		if channel != mockRows[row][2] {
			t.Fatalf("Expecting %v got %v", mockRows[row][2], channel)
		}
		channels = append(channels, channel)
		times = append(times, time)
		addeds = append(addeds, added)
		row++
	}
	if len(times) != len(mockRows) || len(channels) != len(mockRows) || len(addeds) != len(mockRows) {
		t.Error("Did not fetch results properly")
	}
}
