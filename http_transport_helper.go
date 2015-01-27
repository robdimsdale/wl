package wundergo

import (
	"io"
	"net/http"
)

// HTTPTransportHelper provides a boundary object for
// creation and execution of http requests.
type HTTPTransportHelper interface {
	NewRequest(method, urlStr string, body io.Reader) (*http.Request, error)
	DoRequest(req *http.Request) (resp *http.Response, err error)
}

// DefaultHTTPTransportHelper is an implementation of HTTPTransportHelper.
type DefaultHTTPTransportHelper struct {
	client http.Client
}

// NewRequest returns a new http.Request.
func (h DefaultHTTPTransportHelper) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, urlStr, body)
}

// DoRequest executes the provided http.Request, returning its response.
func (h DefaultHTTPTransportHelper) DoRequest(req *http.Request) (resp *http.Response, err error) {
	return h.client.Do(req)
}
