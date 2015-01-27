package wundergo

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

// NewHTTPTransportHelper allows for the injection of a HTTPTransportHelper
// boundary object.
var NewHTTPTransportHelper = func() HTTPTransportHelper {
	return &DefaultHTTPTransportHelper{
		client: http.Client{},
	}
}

// HTTPHelper provides a wrapper around a variety of HTTP methods such as GET
// and POST.
// It is primarily useful to encapsulate things like authentication headers.
type HTTPHelper interface {
	Get(url string) (*http.Response, error)
	Post(url string, body []byte) (*http.Response, error)
	Put(url string, body []byte) (*http.Response, error)
	Patch(url string, body []byte) (*http.Response, error)
	Delete(url string) (*http.Response, error)
}

// OauthClientHTTPHelper implements HTTPHelper, utilizing oath credentials for
// authentication.
type OauthClientHTTPHelper struct {
	accessToken   string
	clientID      string
	httpTransport HTTPTransportHelper
}

// NewOauthClientHTTPHelper provides a conveinient mechanism for initializing an
// OauthClientHTTPHelper.
func NewOauthClientHTTPHelper(accessToken string, clientID string) *OauthClientHTTPHelper {
	return &OauthClientHTTPHelper{
		accessToken:   accessToken,
		clientID:      clientID,
		httpTransport: NewHTTPTransportHelper(),
	}
}

// Get provides a wrapper around http.Get.
// Response is guaranteed to be non-nil if error is nil
func (h OauthClientHTTPHelper) Get(url string) (*http.Response, error) {
	return h.performHTTPAction(
		url,
		"GET",
		nil,
		nil,
	)
}

// Put provides a wrapper around http.Put.
// Response is guaranteed to be non-nil if error is nil
func (h OauthClientHTTPHelper) Put(url string, body []byte) (*http.Response, error) {
	return h.performHTTPAction(
		url,
		"PUT",
		body,
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	)
}

// Post provides a wrapper around http.Post.
// Response is guaranteed to be non-nil if error is nil
func (h OauthClientHTTPHelper) Post(url string, body []byte) (*http.Response, error) {
	return h.performHTTPAction(
		url,
		"POST",
		body,
		map[string]string{"Content-Type": "application/json"})
}

// Patch provides a wrapper around http.Patch.
// Response is guaranteed to be non-nil if error is nil
func (h OauthClientHTTPHelper) Patch(url string, body []byte) (*http.Response, error) {
	return h.performHTTPAction(
		url,
		"PATCH",
		body,
		map[string]string{"Content-Type": "application/json"})
}

// Delete provides a wrapper around http.Delete.
// Response is guaranteed to be non-nil if error is nil
func (h OauthClientHTTPHelper) Delete(url string) (*http.Response, error) {
	return h.performHTTPAction(
		url,
		"DELETE",
		nil,
		nil)
}

func (h OauthClientHTTPHelper) performHTTPAction(
	url string,
	action string,
	body []byte,
	headers map[string]string,
) (*http.Response, error) {

	req, err := h.httpTransport.NewRequest(action, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Access-Token", h.accessToken)
	req.Header.Add("X-Client-ID", h.clientID)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if body != nil && len(body) != 0 {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}

	resp, err := h.httpTransport.DoRequest(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("Nil response returned")
	}

	return resp, nil
}
