// Pacakge nface handles all communication between Go code and the Reddit api.
package nface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ReqAction int
const (
	GET  = iota
	POST = iota
)

const (
	// contentType is a header flag for POST requests so the reddit api
	// knows how to read the request body.
	contentType = "application/x-www-form-urlencoded"
)

type Client struct {
	// client holds an http.Transport that automatically handles OAuth.
	client *http.Client
	// user is a string attached to all request headers that describes
	// the program to the reddit API. (user-agent)
	user string
}

// Request describes how to build an http.Request for the reddit api.
type Request struct {
	// Action is the request type (e.g. "POST" or "GET").
	Action ReqAction
	// BaseUrl is the url of the api call, which values will be appended to.
	BaseUrl string
	// Values holds any parameters for the api call; encoded in url.
	Values *url.Values
}

// NewClient returns a new Client struct.
func NewClient(client *http.Client, userAgent string) (*Client) {
	return &Client{client: client, user: userAgent}
}

// Do executes a request using Client's auth and user agent. The result is
// Unmarshal()ed into response.
func (c *Client) Do(r *Request, response interface{}) error {
	req, err := c.buildRequest(r)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(resp, response)
}

// buildRequest builds an http.Request from a Request struct.
func (c *Client) buildRequest(r *Request) (*http.Request, error) {
	var req *http.Request
	var err error
	if r.Action == GET {
		req, err = getRequest(r.BaseUrl, r.Values)
	} else if r.Action == POST {
		req, err = postRequest(r.BaseUrl, r.Values)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Add("user-agent", c.user)

	return req, nil
}

// doRequest sends a request to the servers and returns the body of the response
// a byte slice.
func (c *Client) doRequest(r *http.Request) ([]byte, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %v\n", resp.StatusCode)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %v", err)
	}

	return buf, nil
}

// postRequest returns a template http.Request with the given url and POST form
// values set.
func postRequest(url string, vals *url.Values) (*http.Request, error) {
	if vals == nil {
		return nil, fmt.Errorf("no values for POST body")
	}

	reqBody := bytes.NewBufferString(vals.Encode())
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", contentType)
	return req, nil
}

// getRequest returns a template http.Request with the given url and GET form
// values set.
func getRequest(url string, vals *url.Values) (*http.Request, error) {
	reqURL := url
	if vals != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, vals.Encode())
	}
	return http.NewRequest("GET", reqURL, nil)
}
