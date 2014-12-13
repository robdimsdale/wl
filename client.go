package wundergo

import (
	"encoding/json"
	"fmt"
)

const (
	apiUrl = "https://a.wunderlist.com/api/v1"
)

var NewLogger = func() Logger {
	return newPrintlnLogger()
}

var NewHTTPHelper = func(accessToken string, clientID string) HTTPHelper {
	return newOauthClientHTTPHelper(accessToken, clientID)
}

type Client interface {
	User() (User, error)
	UpdateUser(user User) (User, error)
	Users() ([]User, error)
	UsersForListID(listId uint) ([]User, error)
}

type OauthClient struct {
	httpHelper HTTPHelper
	logger     Logger
}

func NewOauthClient(accessToken string, clientID string) *OauthClient {
	return &OauthClient{
		httpHelper: NewHTTPHelper(accessToken, clientID),
		logger:     NewLogger(),
	}
}

func (c OauthClient) User() (User, error) {
	b, err := c.httpHelper.Get(fmt.Sprintf("%s/user", apiUrl))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return User{}, err
	}

	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return User{}, err
	}
	return u, nil
}

func (c OauthClient) UpdateUser(user User) (User, error) {
	body := fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name)
	b, err := c.httpHelper.Put(fmt.Sprintf("%s/user", apiUrl), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return User{}, err
	}

	var u User
	err = json.Unmarshal(b, &u)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return User{}, err
	}
	return u, nil
}

func (c OauthClient) Users() ([]User, error) {
	return c.UsersForListID(0)
}

func (c OauthClient) UsersForListID(listId uint) ([]User, error) {
	var b []byte
	var err error

	if listId > 0 {
		b, err = c.httpHelper.Get(fmt.Sprintf("%s/users?list_id=%d", apiUrl, listId))
	} else {
		b, err = c.httpHelper.Get(fmt.Sprintf("%s/users", apiUrl))
	}

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return []User{}, err
	}

	var u []User
	err = json.Unmarshal(b, &u)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return []User{}, err
	}
	return u, nil
}
