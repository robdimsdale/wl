package wundergo

import (
	"encoding/json"
	"fmt"
)

const (
	apiUrl = "https://a.wunderlist.com/api/v1"
)

var NewHTTPHelper = func(accessToken string, clientID string) HTTPHelper {
	return newOauthClientHTTPHelper(accessToken, clientID)
}

type Client interface {
	User() (User, error)
	UpdateUser(user User) (User, error)
	Users() ([]User, error)
}

type OauthClient struct {
	httpHelper HTTPHelper
}

func NewOauthClient(accessToken string, clientID string) *OauthClient {
	return &OauthClient{
		httpHelper: NewHTTPHelper(accessToken, clientID),
	}
}

func (c OauthClient) User() (User, error) {
	b, err := c.httpHelper.Get(fmt.Sprintf("%s/user", apiUrl))
	if err != nil {
		return User{}, err
	}

	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (c OauthClient) UpdateUser(user User) (User, error) {
	body := fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name)
	b, err := c.httpHelper.Put(fmt.Sprintf("%s/user", apiUrl), body)
	if err != nil {
		return User{}, err
	}

	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (c OauthClient) Users() ([]User, error) {
	b, err := c.httpHelper.Get(fmt.Sprintf("%s/users", apiUrl))
	if err != nil {
		return []User{}, err
	}

	var u []User
	err = json.Unmarshal(b, &u)
	if err != nil {
		return []User{}, err
	}
	return u, nil
}
