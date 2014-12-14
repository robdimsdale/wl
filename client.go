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
	return NewOauthClientHTTPHelper(accessToken, clientID)
}

type Client interface {
	User() (User, error)
	UpdateUser(user User) (User, error)
	Users() ([]User, error)
	UsersForListID(listId uint) ([]User, error)
	Lists() ([]List, error)
	List(listID uint) (List, error)
	ListTaskCount(listID uint) (ListTaskCount, error)
	CreateList(listTitle string) (List, error)
	UpdateList(list List) (List, error)
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
	body := []byte(fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name))
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

func (c OauthClient) Lists() ([]List, error) {
	b, err := c.httpHelper.Get(fmt.Sprintf("%s/lists", apiUrl))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return []List{}, err
	}

	var l []List
	err = json.Unmarshal(b, &l)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return []List{}, err
	}
	return l, nil
}

func (c OauthClient) List(listID uint) (List, error) {
	b, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/%d", apiUrl, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return List{}, err
	}

	var l List
	err = json.Unmarshal(b, &l)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return List{}, err
	}
	return l, nil
}

func (c OauthClient) ListTaskCount(listID uint) (ListTaskCount, error) {
	b, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/tasks_count?list_id=%d", apiUrl, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return ListTaskCount{}, err
	}

	var l ListTaskCount
	err = json.Unmarshal(b, &l)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return ListTaskCount{}, err
	}
	return l, nil
}

func (c OauthClient) CreateList(listTitle string) (List, error) {
	body := []byte(fmt.Sprintf(`{"title":"%s"}`, listTitle))
	c.logger.LogLine(fmt.Sprintf("request body: %s", string(body)))
	b, err := c.httpHelper.Post(fmt.Sprintf("%s/lists", apiUrl), body)
	c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return List{}, err
	}

	var l List
	err = json.Unmarshal(b, &l)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return List{}, err
	}
	return l, nil
}

func (c OauthClient) UpdateList(list List) (List, error) {
	body, err := json.Marshal(list)
	if err != nil {

	}
	b, err := c.httpHelper.Patch(fmt.Sprintf("%s/lists/%d", apiUrl, list.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return List{}, err
	}

	var l List
	err = json.Unmarshal(b, &l)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response body: %s", string(b)))
		return List{}, err
	}
	return l, nil
}
