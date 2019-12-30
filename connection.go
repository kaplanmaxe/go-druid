package druid

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

var (
	// ErrPinging is an error returned when health check endpoint returns a non 200.
	ErrPinging = errors.New("druid: error fetching health info from druid")
)

type connection struct {
	Client *http.Client
	Cfg    *Config
}

type queryRequest struct {
	Query        string `json:"query"`
	ResultFormat string `json:"resultFormat"`
	Header       bool   `json:"header"`
}

type queryResponse [][]interface{}

// Prepare implements db.Conn.Prepare and returns a noop statement
func (c *connection) Prepare(stmt string) (driver.Stmt, error) {
	return &statementNoop{}, nil
}

// Close closes a connection.
// TODO: implement
func (c *connection) Close() (err error) {
	return
}

// Begin implements db.Conn.Prepare and is a noop
func (c *connection) Begin() (tx driver.Tx, err error) {
	tx = &transactionNoop{}
	return
}

// Ping implmements db.conn.Prepare and hits the health endpoint of a broker
func (c *connection) Ping(ctx context.Context) (err error) {
	res, err := c.Client.Get(fmt.Sprintf("%s%s", c.Cfg.BrokerAddr, c.Cfg.PingEndpoint))
	if err != nil {
		err = ErrPinging
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("druid: got %d status code from %s", res.StatusCode, c.Cfg.PingEndpoint)
	}
	return
}

// Query queries the druid sql api
// TODO: fix error handling
func (c *connection) Query(q string, args []driver.Value) (driver.Rows, error) {
	return c.query(q, args)
}

func (c *connection) query(q string, args []driver.Value) (*rows, error) {
	queryURL := fmt.Sprintf("%s%s", c.Cfg.BrokerAddr, c.Cfg.QueryEndpoint)
	request := &queryRequest{
		Query:        q,
		ResultFormat: "array",
		Header:       true,
	}
	req, err := json.Marshal(request)
	if err != nil {
		return &rows{}, errors.New("druid: Error marshalling query")
	}
	res, err := c.Client.Post(queryURL, "application/json", bytes.NewReader(req))
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &rows{}, err
	}
	var results queryResponse
	json.Unmarshal(body, &results)
	// No results returned
	// TODO: pass back error from api
	if len(results) == 0 {
		return nil, errors.New("druid: no results returned")
	}
	var columnNames []string

	for _, val := range results[0] {
		// val := reflect.ValueOf(val).Convert(reflect.TypeOf(val))
		columnNames = append(columnNames, val.(string))
	}
	var returnedRows [][]field
	for i := 1; i < len(results); i++ {
		var cols []field
		for _, val := range results[i] {
			cols = append(cols, field{Value: reflect.ValueOf(val), Type: reflect.TypeOf(val)})
		}
		returnedRows = append(returnedRows, cols)
	}

	resultSet := resultSet{
		columnNames: columnNames,
		rows:        returnedRows,
		currentRow:  0,
	}
	r := &rows{
		conn:      c,
		resultSet: resultSet,
	}
	return r, nil
}
