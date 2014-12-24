package wundergo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiUrl = "https://a.wunderlist.com/api/v1"
)

var NewLogger = func() Logger {
	return NewPrintlnLogger()
}

var NewHTTPHelper = func(accessToken string, clientID string) HTTPHelper {
	return NewOauthClientHTTPHelper(accessToken, clientID)
}

var NewJSONHelper = func() JSONHelper {
	return NewDefaultJSONHelper()
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
	DeleteList(list List) error
}

type OauthClient struct {
	httpHelper HTTPHelper
	logger     Logger
	jsonHelper JSONHelper
}

func NewOauthClient(accessToken string, clientID string) *OauthClient {
	return &OauthClient{
		httpHelper: NewHTTPHelper(accessToken, clientID),
		logger:     NewLogger(),
		jsonHelper: NewJSONHelper(),
	}
}

func (c OauthClient) User() (User, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/user", apiUrl))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return User{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return User{}, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &User{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return User{}, err
	}
	return *(u.(*User)), nil
}

func (c OauthClient) readResponseBody(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, errors.New(fmt.Sprintf("Nil body on http response: %v", resp))
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c OauthClient) UpdateUser(user User) (User, error) {
	body := []byte(fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name))
	resp, err := c.httpHelper.Put(fmt.Sprintf("%s/user", apiUrl), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return User{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return User{}, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &User{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return User{}, err
	}
	return *(u.(*User)), nil
}

func (c OauthClient) Users() ([]User, error) {
	return c.UsersForListID(0)
}

func (c OauthClient) UsersForListID(listId uint) ([]User, error) {
	var resp *http.Response
	var err error

	if listId > 0 {
		resp, err = c.httpHelper.Get(fmt.Sprintf("%s/users?list_id=%d", apiUrl, listId))
	} else {
		resp, err = c.httpHelper.Get(fmt.Sprintf("%s/users", apiUrl))
	}

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return []User{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return []User{}, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &([]User{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return []User{}, err
	}
	return *(u.(*[]User)), nil
}

func (c OauthClient) Lists() ([]List, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists", apiUrl))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return []List{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return []List{}, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &([]List{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return []List{}, err
	}
	return *(l.(*[]List)), nil
}

func (c OauthClient) List(listID uint) (List, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/%d", apiUrl, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &List{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}
	return *(l.(*List)), nil
}

func (c OauthClient) ListTaskCount(listID uint) (ListTaskCount, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/tasks_count?list_id=%d", apiUrl, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return ListTaskCount{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return ListTaskCount{}, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &ListTaskCount{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return ListTaskCount{}, err
	}
	return *(l.(*ListTaskCount)), nil
}

func (c OauthClient) CreateList(listTitle string) (List, error) {
	body := []byte(fmt.Sprintf(`{"title":"%s"}`, listTitle))

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/lists", apiUrl), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &List{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}
	return *(l.(*List)), nil
}

func (c OauthClient) UpdateList(list List) (List, error) {
	body, err := c.jsonHelper.Marshal(list)
	if err != nil {
		return List{}, err
	}
	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/lists/%d", apiUrl, list.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &List{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return List{}, err
	}
	return *(l.(*List)), nil
}

func (c OauthClient) DeleteList(list List) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/lists/%d?revision=%d", apiUrl, list.ID, list.Revision))

	if err != nil {
		return err
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}
