package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pivotal-golang/lager"
	"github.com/robdimsdale/wundergo"
)

// User returns the currently logged in user.
// This makes it a good method to validate the auth credentials provided
// in NewoauthClient.
func (c oauthClient) User() (wundergo.User, error) {
	url := fmt.Sprintf(
		"%s/user",
		c.apiURL,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.User{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.User{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.User{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	if resp.Body == nil {
		return wundergo.User{}, errors.New("Nil body returned")
	}

	user := wundergo.User{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.User{}, err
	}

	return user, nil
}

// UpdateUser is a currently undocumented method which updates the provided user.
// Currently the only field that is updated is user.Name
func (c oauthClient) UpdateUser(user wundergo.User) (wundergo.User, error) {
	body := []byte(fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name))
	url := fmt.Sprintf("%s/user", c.apiURL)

	req, err := c.newPutRequest(url, body)
	if err != nil {
		return wundergo.User{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.User{}, err
	}

	if resp.Body == nil {
		return wundergo.User{}, errors.New("Nil body returned")
	}

	returnedUser := wundergo.User{}
	err = json.NewDecoder(resp.Body).Decode(&returnedUser)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.User{}, err
	}

	return returnedUser, nil
}

// Users returns a list of all users the client can access.
func (c oauthClient) Users() ([]wundergo.User, error) {
	return c.UsersForListID(0)
}

// UsersForListID returns a list of users the client can access,
// restricted to users that have access to the provided list.
func (c oauthClient) UsersForListID(listID uint) ([]wundergo.User, error) {
	var url string
	if listID > 0 {
		url = fmt.Sprintf("%s/users?list_id=%d", c.apiURL, listID)
	} else {
		url = fmt.Sprintf("%s/users", c.apiURL)
	}

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	if resp.Body == nil {
		return nil, errors.New("Nil body returned")
	}

	users := []wundergo.User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return nil, err
	}

	return users, nil
}
