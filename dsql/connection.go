package dsql

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
)

var (
	// ErrPinging is an error returned when health check endpoint returns a non 200.
	ErrPinging = errors.New("druid: error fetching health info from druid")
	// ErrCancelled is an error returned when we receive a cancellation event from a context object
	ErrCancelled = errors.New("druid: cancellation received")
	// ErrRequestForm is an error returned when failing to form a request
	ErrRequestForm = errors.New("druid: error forming request")
)

type key int

const (
	transportKey key = iota
	requestKey
)

type connection struct {
	Client    *http.Client
	Cfg       *Config
	closeCh   chan struct{}
	watcherCh chan context.Context
	errorCh   chan error
	resultsCh chan []byte
	requestCh chan *http.Request
	closed    bool
	mtx       sync.Mutex
}

type queryRequest struct {
	Query        string `json:"query"`
	ResultFormat string `json:"resultFormat"`
	Header       bool   `json:"header"`
}

type queryResponse [][]interface{}

// Prepare implements db.Conn.Prepare and returns a noop statement
func (c *connection) Prepare(stmt string) (driver.Stmt, error) {
	return &statementNoop{}, driver.ErrSkip
}

// Close closes a connection.
func (c *connection) Close() (err error) {
	if c.closed == true {
		return
	}
	c.mtx.Lock()
	c.closed = true
	c.mtx.Unlock()
	close(c.closeCh)
	c.Client = nil
	return
}

// Begin implements db.Conn.Prepare and is a noop
func (c *connection) Begin() (tx driver.Tx, err error) {
	tx = &transactionNoop{}
	return tx, driver.ErrSkip
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
		err = ErrPinging
	}
	return
}

func (c *connection) startRequestPipeline() {
	go func() {
		for {
			select {
			case req := <-c.requestCh:
				res, err := c.Client.Do(req)
				if err != nil {
					c.errorCh <- err
					return
				}
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					c.errorCh <- err
				}
				c.resultsCh <- body
			case <-c.closeCh:
				return
			}
		}
	}()
}

// Query queries the druid sql api
func (c *connection) Query(q string, args []driver.Value) (driver.Rows, error) {
	return c.query(q, args)
}

func (c *connection) makeRequest(q string) (req *http.Request, err error) {
	queryURL := fmt.Sprintf("%s%s", c.Cfg.BrokerAddr, c.Cfg.QueryEndpoint)
	request := &queryRequest{
		Query:        q,
		ResultFormat: "array",
		Header:       true,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, ErrRequestForm
	}
	req, err = http.NewRequest("POST", queryURL, bytes.NewReader(payload))
	if err != nil {
		return nil, ErrRequestForm
	}
	return
}

func (c *connection) parseResponse(body []byte) (r *rows, err error) {
	var results queryResponse
	err = json.Unmarshal(body, &results)
	if err != nil {
		return &rows{}, err
	}
	// No results returned
	if len(results) == 0 {
		return &rows{}, sql.ErrNoRows
	}
	var columnNames []string

	for _, val := range results[0] {
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
	r = &rows{
		conn:      c,
		resultSet: resultSet,
	}
	return r, nil
}

// TODO: better error hanlding
func (c *connection) query(q string, args []driver.Value) (*rows, error) {
	req, err := c.makeRequest(q)
	if err != nil {
		return &rows{}, err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return &rows{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &rows{}, err
	}
	return c.parseResponse(body)
}

func (c *connection) queryContext(ctx context.Context, q string, args []driver.NamedValue) (*rows, error) {
	req, err := c.makeRequest(q)
	if err != nil {
		return &rows{}, errors.New("druid: error making request")
	}
	req.WithContext(ctx)
	c.requestCh <- req
	tr := &http.Transport{}
	c.Client.Transport = tr
	var r *rows
	select {
	case body := <-c.resultsCh:
		r, err = c.parseResponse(body)
		if err != nil {
			return r, err
		}
	case err = <-c.errorCh:
	case <-ctx.Done():
		tr.CancelRequest(req)
		err = ctx.Err()
		return r, err
	}
	return r, err
}
func (c *connection) QueryContext(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	vals, err := namedValuesToValues(args)
	if err != nil {
		return nil, err
	}

	if ctx.Done() == nil {
		return c.query(q, vals)
	}
	return c.queryContext(ctx, q, args)
}
