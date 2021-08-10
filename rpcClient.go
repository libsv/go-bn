package bitcoin

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	rpcClientTimeout = 120
)

// A rpcClient represents a JSON RPC client (over HTTP(s)).
type rpcClient struct {
	serverAddr string
	user       string
	passwd     string
	httpClient *http.Client
}

// rpcRequest represent a RCP request
type rpcRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
	JSONRpc string      `json:"jsonrpc"`
}

// rpcError represents a RCP error
/*type rpcError struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}*/

// rpcResponse represents a RCP response
type rpcResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    interface{}     `json:"error"`
}

func newClient(host string, port int, user, passwd string, useSSL bool) (c *rpcClient, err error) {
	if len(host) == 0 {
		err = errors.New("Bad call missing argument host")
		return
	}
	var serverAddr string
	var httpClient *http.Client
	if useSSL {
		serverAddr = "https://"
		t := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: t}
	} else {
		serverAddr = "http://"
		httpClient = &http.Client{}
	}
	c = &rpcClient{serverAddr: fmt.Sprintf("%s%s:%d", serverAddr, host, port), user: user, passwd: passwd, httpClient: httpClient}
	return
}

// doTimeoutRequest process a HTTP request with timeout
func (c *rpcClient) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("Timeout reading data from server")
	}
}

// call prepare & exec the request
func (c *rpcClient) call(method string, params interface{}) (rpcResponse, error) {
	connectTimer := time.NewTimer(rpcClientTimeout * time.Second)
	rpcR := rpcRequest{method, params, time.Now().UnixNano(), "1.0"}
	payloadBuffer := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(payloadBuffer)

	err := jsonEncoder.Encode(rpcR)
	if err != nil {
		return rpcResponse{}, fmt.Errorf("failed to encode rpc request: %w", err)
	}

	req, err := http.NewRequest("POST", c.serverAddr, payloadBuffer)
	if err != nil {
		return rpcResponse{}, fmt.Errorf("failed to create new http request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	// Auth ?
	if len(c.user) > 0 || len(c.passwd) > 0 {
		req.SetBasicAuth(c.user, c.passwd)
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return rpcResponse{}, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rpcResponse{}, fmt.Errorf("failed to read response: %w", err)
	}

	var rr rpcResponse

	if resp.StatusCode != 200 {
		json.Unmarshal(data, &rr)
		v, ok := rr.Err.(map[string]interface{})
		if ok {
			err = errors.New(v["message"].(string))
		} else {
			err = errors.New("HTTP error: " + resp.Status)
		}

		return rr, fmt.Errorf("unexpected response code 200: %w", err)
	}

	err = json.Unmarshal(data, &rr)
	if err != nil {
		return rr, fmt.Errorf("failed to unmarshall response: %w", err)
	}

	return rr, nil
}
