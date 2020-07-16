package axcessms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/context/ctxhttp"
)

// Client represents a full API client instance.
type Client struct {
	conn *http.Client

	DebugWriter io.Writer

	testMode bool

	tokenProvider TokenProvider
}

const (
	// TestHost is the base API URL to be used for test-mode interactions
	TestHost = "https://test.oppwa.com"

	// LiveHost is the base API URL to be used for live-mode interactions
	LiveHost = "https://oppwa.com"
)

// New creates a new client
func New(ctx context.Context, provider TokenProvider) *Client {
	hc := &http.Client{
		Transport: http.DefaultTransport,
	}

	return &Client{
		conn:          hc,
		testMode:      false,
		tokenProvider: provider,
	}
}

// Do performs the given HTTP request on the client
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Axcessms-go +https://github.com/mrzen/axcessms")

	token, err := c.tokenProvider()

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	if c.testMode && c.DebugWriter != nil {
		if req.Body != nil {
			fmt.Fprintln(c.DebugWriter, "---- REQUEST ----")

			cb := new(bytes.Buffer)
			tr := io.TeeReader(req.Body, cb)
			req.Body = ioutil.NopCloser(tr)
			req.Write(c.DebugWriter)
			req.Body = ioutil.NopCloser(bytes.NewReader(cb.Bytes()))
		}
	}

	return ctxhttp.Do(ctx, c.conn, req)
}

// URLEncodable is a type which is URL Encodable
type URLEncodable interface {
	URLEncode() url.Values
}

// Run sends the given request body as a WWW-Form to the given path, and decodes a JSON response
// into the given interface, returning any error
func (c *Client) Run(ctx context.Context, method, path string, body URLEncodable, into interface{}) error {

	var reader io.ReadCloser

	if body != nil {
		reader = ioutil.NopCloser(bytes.NewBufferString(body.URLEncode().Encode()))
	}

	req, err := http.NewRequest(method, c.getEndpoint()+path, reader)
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	resp, err := c.Do(ctx, req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c.testMode && c.DebugWriter != nil {
		cb := new(bytes.Buffer)
		tr := io.TeeReader(resp.Body, cb)

		resp.Body = ioutil.NopCloser(tr)

		fmt.Fprintln(c.DebugWriter, "--- Response")

		resp.Write(c.DebugWriter)

		resp.Body = ioutil.NopCloser(bytes.NewReader(cb.Bytes()))
	}

	if resp.StatusCode >= 400 {
		errb := APIResponse{}

		cb := new(bytes.Buffer)
		tr := io.TeeReader(resp.Body, cb)

		if err := json.NewDecoder(tr).Decode(&errb); err != nil {
			return err
		}

		return errors.New(errb.Result.Description)
	}

	return json.NewDecoder(resp.Body).Decode(&into)

}

// SetTestMode sets if the API will communicate in test or production mode
func (c *Client) SetTestMode(test bool) {
	c.testMode = test
}

func (c Client) getEndpoint() string {
	if c.testMode {
		return TestHost
	}

	return LiveHost
}
