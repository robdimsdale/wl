package wundergo

import (
	"io"
	"net/http"
)

type HTTPTransportHelper interface {
	NewRequest(method, urlStr string, body io.Reader) (*http.Request, error)
	DoRequest(req *http.Request) (resp *http.Response, err error)
}

type DefaultHTTPTransportHelper struct {
	client http.Client
}

func NewDefaultHTTPTransportHelper() *DefaultHTTPTransportHelper {
	return &DefaultHTTPTransportHelper{
		client: http.Client{},
	}
}

func (h DefaultHTTPTransportHelper) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, urlStr, body)
}

func (h DefaultHTTPTransportHelper) DoRequest(req *http.Request) (resp *http.Response, err error) {
	return h.client.Do(req)
}
