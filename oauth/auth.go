package oauth

import (
	"fmt"
	"net/http"
)

// Authed returns (true, nil) if the user is successfully logged-in.
// It returns (false, nil) if the user is not logged-in, and
// (false, err) if there is an error determining whether the user is logged-in.
// The HTTP status codes 401 and 403 are considered not logged-in, rather than
// errors; all other status codes greater than 400 are considered errors.
func (c oauthClient) Authed() (bool, error) {
	// /user is a good endpoint to determine whether a user is logged-in.
	url := fmt.Sprintf(
		"%s/user",
		c.apiURL,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.do(req)
	if err != nil {
		return false, err
	}

	// For auth purposes we only care about the status code, not the contents of
	// the body.
	// We trust that the status code is accurate and sufficient
	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return false, nil
	default:
		return false, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}
}
