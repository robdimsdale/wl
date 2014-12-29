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
	User() (*User, error)
	UpdateUser(user User) (*User, error)
	Users() (*[]User, error)
	UsersForListID(listID uint) (*[]User, error)
	Lists() (*[]List, error)
	List(listID uint) (*List, error)
	ListTaskCount(listID uint) (*ListTaskCount, error)
	CreateList(title string) (*List, error)
	UpdateList(list List) (*List, error)
	DeleteList(list List) error
	NotesForListID(listID uint) (*[]Note, error)
	NotesForTaskID(taskID uint) (*[]Note, error)
	Note(noteID uint) (*Note, error)
	CreateNote(content string, taskID uint) (*Note, error)
	UpdateNote(note Note) (*Note, error)
	DeleteNote(note Note) error
	TasksForListID(listID uint) (*[]Task, error)
	CompletedTasksForListID(listID uint, completed bool) (*[]Task, error)
	Task(taskID uint) (*Task, error)
	CreateTask(
		title string,
		listID uint,
		assigneeID uint,
		completed bool,
		recurrenceType string,
		recurrenceCount uint,
		dueDate string,
		starred bool,
	) (*Task, error)
	UpdateTask(task Task) (*Task, error)
	DeleteTask(task Task) error
	SubtasksForListID(listID uint) (*[]Subtask, error)
	SubtasksForTaskID(taskID uint) (*[]Subtask, error)
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

func (c OauthClient) User() (*User, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/user", apiUrl))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &User{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return u.(*User), nil
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

func (c OauthClient) UpdateUser(user User) (*User, error) {
	body := []byte(fmt.Sprintf("revision=%d&name=%s", user.Revision, user.Name))
	resp, err := c.httpHelper.Put(fmt.Sprintf("%s/user", apiUrl), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &User{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return u.(*User), nil
}

func (c OauthClient) Users() (*[]User, error) {
	return c.UsersForListID(0)
}

func (c OauthClient) UsersForListID(listID uint) (*[]User, error) {
	var resp *http.Response
	var err error

	if listID > 0 {
		resp, err = c.httpHelper.Get(fmt.Sprintf("%s/users?list_id=%d", apiUrl, listID))
	} else {
		resp, err = c.httpHelper.Get(fmt.Sprintf("%s/users", apiUrl))
	}

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &([]User{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return u.(*[]User), nil
}

func (c OauthClient) Lists() (*[]List, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists", apiUrl))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &([]List{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*[]List), nil
}

func (c OauthClient) List(listID uint) (*List, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/%d", apiUrl, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &List{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*List), nil
}

func (c OauthClient) ListTaskCount(listID uint) (*ListTaskCount, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/tasks_count?list_id=%d", apiUrl, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &ListTaskCount{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*ListTaskCount), nil
}

func (c OauthClient) CreateList(title string) (*List, error) {
	body := []byte(fmt.Sprintf(`{"title":"%s"}`, title))

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/lists", apiUrl), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &List{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*List), nil
}

func (c OauthClient) UpdateList(list List) (*List, error) {
	body, err := c.jsonHelper.Marshal(list)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/lists/%d", apiUrl, list.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &List{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*List), nil
}

func (c OauthClient) DeleteList(list List) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/lists/%d?revision=%d", apiUrl, list.ID, list.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent))
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) NotesForListID(listID uint) (*[]Note, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/notes?list_id=%d", apiUrl, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	n, err := c.jsonHelper.Unmarshal(b, &[]Note{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return n.(*[]Note), nil
}

func (c OauthClient) NotesForTaskID(taskID uint) (*[]Note, error) {
	var resp *http.Response
	var err error

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/notes?task_id=%d", apiUrl, taskID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	n, err := c.jsonHelper.Unmarshal(b, &[]Note{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return n.(*[]Note), nil
}

func (c OauthClient) Note(noteID uint) (*Note, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/notes/%d", apiUrl, noteID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	n, err := c.jsonHelper.Unmarshal(b, &Note{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return n.(*Note), nil
}

func (c OauthClient) CreateNote(content string, taskID uint) (*Note, error) {
	body := []byte(fmt.Sprintf(`{"content":"%s","task_id":%d}`, content, taskID))

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/notes", apiUrl), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	n, err := c.jsonHelper.Unmarshal(b, &Note{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return n.(*Note), nil
}

func (c OauthClient) UpdateNote(note Note) (*Note, error) {
	body, err := c.jsonHelper.Marshal(note)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/notes/%d", apiUrl, note.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	n, err := c.jsonHelper.Unmarshal(b, &Note{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return n.(*Note), nil
}

func (c OauthClient) DeleteNote(note Note) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/notes/%d?revision=%d", apiUrl, note.ID, note.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent))
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) TasksForListID(listID uint) (*[]Task, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/tasks?list_id=%d", apiUrl, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	t, err := c.jsonHelper.Unmarshal(b, &[]Task{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return t.(*[]Task), nil
}

func (c OauthClient) CompletedTasksForListID(listID uint, completed bool) (*[]Task, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/tasks?list_id=%d&completed=%t", apiUrl, listID, completed))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	t, err := c.jsonHelper.Unmarshal(b, &[]Task{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return t.(*[]Task), nil
}

func (c OauthClient) Task(taskID uint) (*Task, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/tasks/%d", apiUrl, taskID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	t, err := c.jsonHelper.Unmarshal(b, &Task{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return t.(*Task), nil
}

type taskCreateConfig struct {
	ListID          uint   `json:"list_id"`
	Title           string `json:"title"`
	AssigneeID      uint   `json:"assignee_id,omitempty"`
	Completed       bool   `json:"completed,omitempty"`
	RecurrenceType  string `json:"recurrence_type,omitempty"`
	RecurrenceCount uint   `json:"recurrence_count,omitempty"`
	DueDate         string `json:"due_date,omitempty"`
	Starred         bool   `json:"starred,omitempty"`
}

type TaskUpdateConfig struct {
	Title           string   `json:"title"`
	Revision        uint     `json:"revision"`
	AssigneeID      uint     `json:"assignee_id,omitempty"`
	Completed       bool     `json:"completed,omitempty"`
	RecurrenceType  string   `json:"recurrence_type,omitempty"`
	RecurrenceCount uint     `json:"recurrence_count,omitempty"`
	DueDate         string   `json:"due_date,omitempty"`
	Starred         bool     `json:"starred,omitempty"`
	Remove          []string `json:"remove,omitempty"`
}

func (c OauthClient) validateRecurrence(recurrenceType string, recurrenceCount uint) error {
	if recurrenceType == "" && recurrenceCount > 0 {
		return errors.New("recurrenceCount must be zero if provided recurrenceType is not provided")
	}

	if recurrenceCount == 0 && recurrenceType != "" {
		return errors.New("recurrenceType must be valid if provided recurrenceCount is non-zero")
	}

	return nil
}

func (c OauthClient) CreateTask(
	title string,
	listID uint,
	assigneeID uint,
	completed bool,
	recurrenceType string,
	recurrenceCount uint,
	dueDate string,
	starred bool,
) (*Task, error) {

	err := c.validateRecurrence(recurrenceType, recurrenceCount)
	if err != nil {
		return nil, err
	}

	tcc := taskCreateConfig{
		ListID:          listID,
		Title:           title,
		AssigneeID:      assigneeID,
		Completed:       completed,
		RecurrenceType:  recurrenceType,
		RecurrenceCount: recurrenceCount,
		DueDate:         dueDate,
		Starred:         starred,
	}

	body, err := c.jsonHelper.Marshal(tcc)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/tasks", apiUrl), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	t, err := c.jsonHelper.Unmarshal(b, &Task{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return t.(*Task), nil
}

func (c OauthClient) UpdateTask(task Task) (*Task, error) {
	err := c.validateRecurrence(task.RecurrenceType, task.RecurrenceCount)
	if err != nil {
		return nil, err
	}

	origTask, err := c.Task(task.ID)
	if err != nil {
		return nil, err
	}

	tuc := TaskUpdateConfig{
		Title:    task.Title,
		Revision: task.Revision,
		Remove:   []string{},
	}

	if origTask.AssigneeID == task.AssigneeID {
		tuc.AssigneeID = origTask.AssigneeID
	} else {
		if task.AssigneeID == 0 {
			tuc.Remove = append(tuc.Remove, "assignee_id")
		} else {
			tuc.AssigneeID = task.AssigneeID
		}
	}

	if origTask.DueDate == task.DueDate {
		tuc.DueDate = origTask.DueDate
	} else {
		if task.DueDate == "" {
			tuc.Remove = append(tuc.Remove, "due_date")
		} else {
			tuc.DueDate = task.DueDate
		}
	}

	if origTask.RecurrenceCount == task.RecurrenceCount && origTask.RecurrenceType == task.RecurrenceType {
		tuc.RecurrenceCount = origTask.RecurrenceCount
		tuc.RecurrenceType = origTask.RecurrenceType
	} else {
		if task.RecurrenceCount == 0 {
			tuc.Remove = append(tuc.Remove, "recurrence_type")
			tuc.Remove = append(tuc.Remove, "recurrence_count")
		} else {
			tuc.RecurrenceCount = task.RecurrenceCount
			tuc.RecurrenceType = task.RecurrenceType
		}
	}

	tuc.Completed = task.Completed
	tuc.Starred = task.Starred

	body, err := c.jsonHelper.Marshal(tuc)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/tasks/%d", apiUrl, task.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	t, err := c.jsonHelper.Unmarshal(b, &Task{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return t.(*Task), nil
}

func (c OauthClient) DeleteTask(task Task) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/tasks/%d?revision=%d", apiUrl, task.ID, task.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent))
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) SubtasksForListID(listID uint) (*[]Subtask, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/subtasks?list_id=%d", apiUrl, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	s, err := c.jsonHelper.Unmarshal(b, &[]Subtask{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return s.(*[]Subtask), nil
}

func (c OauthClient) SubtasksForTaskID(taskID uint) (*[]Subtask, error) {
	var resp *http.Response
	var err error

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/subtasks?task_id=%d", apiUrl, taskID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK))
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	s, err := c.jsonHelper.Unmarshal(b, &[]Subtask{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return s.(*[]Subtask), nil
}
