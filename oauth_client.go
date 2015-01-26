package wundergo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewLogger allows for the injection of a Logger boundary object.
// Isolates the client from logging implementation details.
// Defaults to return NewPrintlnLogger.
var NewLogger = func() Logger {
	return NewPrintlnLogger()
}

// NewHTTPHelper allows for the injection of an HTTPHelper boundary object.
// Isolates the client from directly interacting with HTTP - adding auth headers etc.
// Defaults to NewOauthClientHTTPHelper.
var NewHTTPHelper = func(accessToken string, clientID string) HTTPHelper {
	return NewOauthClientHTTPHelper(accessToken, clientID)
}

// NewJSONHelper allows for the injection of a JSONHelper boundary object.
// Isolates the client from JSON marshalling and unmarshalling.
// Defaults to NewDefaultJSONHelper.
var NewJSONHelper = func() JSONHelper {
	return NewDefaultJSONHelper()
}

// OauthClient is an implementation of Client.
type OauthClient struct {
	httpHelper HTTPHelper
	logger     Logger
	jsonHelper JSONHelper
}

// NewOauthClient is a utility method to simplify initialization
// of a new OauthClient.
func NewOauthClient(accessToken string, clientID string) *OauthClient {
	return &OauthClient{
		httpHelper: NewHTTPHelper(accessToken, clientID),
		logger:     NewLogger(),
		jsonHelper: NewJSONHelper(),
	}
}

func (c OauthClient) User() (*User, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/user", apiURL))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
		return nil, fmt.Errorf("Nil body on http response: %v", resp)
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
	resp, err := c.httpHelper.Put(fmt.Sprintf("%s/user", apiURL), body)
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
		resp, err = c.httpHelper.Get(fmt.Sprintf("%s/users?list_id=%d", apiURL, listID))
	} else {
		resp, err = c.httpHelper.Get(fmt.Sprintf("%s/users", apiURL))
	}

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists", apiURL))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/%d", apiURL, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/lists/tasks_count?list_id=%d", apiURL, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/lists", apiURL), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
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

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/lists/%d", apiURL, list.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/lists/%d?revision=%d", apiURL, list.ID, list.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
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

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/notes?list_id=%d", apiURL, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/notes?task_id=%d", apiURL, taskID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/notes/%d", apiURL, noteID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	body := []byte(fmt.Sprintf(`{"content":"%s","task_id":%d}`, content, taskID))

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/notes", apiURL), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
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

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/notes/%d", apiURL, note.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/notes/%d?revision=%d", apiURL, note.ID, note.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
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

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/tasks?list_id=%d", apiURL, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/tasks?list_id=%d&completed=%t", apiURL, listID, completed))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/tasks/%d", apiURL, taskID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

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

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/tasks", apiURL), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
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

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/tasks/%d", apiURL, task.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/tasks/%d?revision=%d", apiURL, task.ID, task.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
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

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/subtasks?list_id=%d", apiURL, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/subtasks?task_id=%d", apiURL, taskID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

func (c OauthClient) CompletedSubtasksForListID(listID uint, completed bool) (*[]Subtask, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/subtasks?list_id=%d&completed=%t", apiURL, listID, completed))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

func (c OauthClient) CompletedSubtasksForTaskID(taskID uint, completed bool) (*[]Subtask, error) {
	var resp *http.Response
	var err error

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/subtasks?task_id=%d&completed=%t", apiURL, taskID, completed))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
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

func (c OauthClient) Subtask(subtaskID uint) (*Subtask, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/subtasks/%d", apiURL, subtaskID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	s, err := c.jsonHelper.Unmarshal(b, &Subtask{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return s.(*Subtask), nil
}

func (c OauthClient) CreateSubtask(
	title string,
	taskID uint,
	completed bool,
) (*Subtask, error) {

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	body := []byte(fmt.Sprintf(`{"title":"%s","task_id":%d,"completed":%t}`, title, taskID, completed))

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/subtasks", apiURL), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	s, err := c.jsonHelper.Unmarshal(b, &Subtask{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return s.(*Subtask), nil
}

func (c OauthClient) UpdateSubtask(subtask Subtask) (*Subtask, error) {
	body, err := c.jsonHelper.Marshal(subtask)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/subtasks/%d", apiURL, subtask.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	s, err := c.jsonHelper.Unmarshal(b, &Subtask{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return s.(*Subtask), nil
}

func (c OauthClient) DeleteSubtask(subtask Subtask) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/subtasks/%d?revision=%d", apiURL, subtask.ID, subtask.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) RemindersForListID(listID uint) (*[]Reminder, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/reminders?list_id=%d", apiURL, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	r, err := c.jsonHelper.Unmarshal(b, &[]Reminder{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return r.(*[]Reminder), nil
}

func (c OauthClient) RemindersForTaskID(taskID uint) (*[]Reminder, error) {
	var resp *http.Response
	var err error

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/reminders?task_id=%d", apiURL, taskID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	r, err := c.jsonHelper.Unmarshal(b, &[]Reminder{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return r.(*[]Reminder), nil
}

func (c OauthClient) Reminder(reminderID uint) (*Reminder, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/reminders/%d", apiURL, reminderID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	r, err := c.jsonHelper.Unmarshal(b, &Reminder{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return r.(*Reminder), nil
}

func (c OauthClient) CreateReminder(
	date string,
	taskID uint,
	createdByDeviceUdid string,
) (*Reminder, error) {

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	var body []byte
	if createdByDeviceUdid == "" {
		body = []byte(fmt.Sprintf(`{"date":"%s","task_id":%d}`, date, taskID))
	} else {
		body = []byte(fmt.Sprintf(`{"date":"%s","task_id":%d,"created_by_device_udid":%s}`, date, taskID, createdByDeviceUdid))
	}

	resp, err := c.httpHelper.Post(fmt.Sprintf("%s/reminders", apiURL), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	r, err := c.jsonHelper.Unmarshal(b, &Reminder{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return r.(*Reminder), nil
}

func (c OauthClient) UpdateReminder(reminder Reminder) (*Reminder, error) {
	body, err := c.jsonHelper.Marshal(reminder)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/reminders/%d", apiURL, reminder.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	r, err := c.jsonHelper.Unmarshal(b, &Reminder{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return r.(*Reminder), nil
}

func (c OauthClient) DeleteReminder(reminder Reminder) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/reminders/%d?revision=%d", apiURL, reminder.ID, reminder.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) ListPositions() (*[]Position, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/list_positions", apiURL))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &([]Position{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*[]Position), nil
}

func (c OauthClient) ListPosition(listPositionID uint) (*Position, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/list_positions/%d", apiURL, listPositionID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &Position{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*Position), nil
}

func (c OauthClient) UpdateListPosition(listPosition Position) (*Position, error) {
	body, err := c.jsonHelper.Marshal(listPosition)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/list_positions/%d", apiURL, listPosition.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	p, err := c.jsonHelper.Unmarshal(b, &Position{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return p.(*Position), nil
}

func (c OauthClient) TaskPositionsForListID(listID uint) (*[]Position, error) {

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/task_positions?list_id=%d", apiURL, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &([]Position{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*[]Position), nil
}

func (c OauthClient) TaskPosition(taskPositionID uint) (*Position, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/task_positions/%d", apiURL, taskPositionID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &Position{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*Position), nil
}

func (c OauthClient) UpdateTaskPosition(taskPosition Position) (*Position, error) {
	body, err := c.jsonHelper.Marshal(taskPosition)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/task_positions/%d", apiURL, taskPosition.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	p, err := c.jsonHelper.Unmarshal(b, &Position{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return p.(*Position), nil
}

func (c OauthClient) SubtaskPositionsForListID(listID uint) (*[]Position, error) {

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/subtask_positions?list_id=%d", apiURL, listID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &([]Position{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*[]Position), nil
}

func (c OauthClient) SubtaskPositionsForTaskID(taskID uint) (*[]Position, error) {

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/subtask_positions?task_id=%d", apiURL, taskID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &([]Position{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*[]Position), nil
}

func (c OauthClient) SubtaskPosition(subTaskPositionID uint) (*Position, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/subtask_positions/%d", apiURL, subTaskPositionID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &Position{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*Position), nil
}

func (c OauthClient) UpdateSubtaskPosition(subTaskPosition Position) (*Position, error) {
	body, err := c.jsonHelper.Marshal(subTaskPosition)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/subtask_positions/%d", apiURL, subTaskPosition.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	p, err := c.jsonHelper.Unmarshal(b, &Position{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return p.(*Position), nil
}

func (c OauthClient) Memberships() (*[]Membership, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/memberships", apiURL))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	l, err := c.jsonHelper.Unmarshal(b, &([]Membership{}))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return l.(*[]Membership), nil
}

func (c OauthClient) Membership(membershipID uint) (*Membership, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/memberships/%d", apiURL, membershipID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	m, err := c.jsonHelper.Unmarshal(b, &Membership{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return m.(*Membership), nil
}

func (c OauthClient) MembershipsForListID(listID uint) (*[]Membership, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/memberships?list_id=%d", apiURL, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	m, err := c.jsonHelper.Unmarshal(b, &[]Membership{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return m.(*[]Membership), nil
}

func (c OauthClient) AddMemberToListViaUserID(userID uint, listID uint, muted bool) (*Membership, error) {
	var resp *http.Response
	var err error

	if userID == 0 {
		return nil, errors.New("userID must be > 0")
	}

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Post(fmt.Sprintf("%s/memberships?user_id=%d&list_id=%d&muted=%t", apiURL, userID, listID, muted), nil)

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	m, err := c.jsonHelper.Unmarshal(b, &Membership{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return m.(*Membership), nil
}

func (c OauthClient) AddMemberToListViaEmailAddress(emailAddress string, listID uint, muted bool) (*Membership, error) {
	var resp *http.Response
	var err error

	if emailAddress == "" {
		return nil, errors.New("emailAddress must not be empty")
	}

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Post(fmt.Sprintf("%s/memberships?email=%s&list_id=%d&muted=%t", apiURL, emailAddress, listID, muted), nil)

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	m, err := c.jsonHelper.Unmarshal(b, &Membership{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return m.(*Membership), nil
}

func (c OauthClient) AcceptMember(membership Membership) (*Membership, error) {
	body, err := c.jsonHelper.Marshal(membership)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/memberships/%d", apiURL, membership.ID), body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	m, err := c.jsonHelper.Unmarshal(b, &Membership{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return m.(*Membership), nil
}

func (c OauthClient) RejectInvite(membership Membership) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/memberships/%d?revision=%d", apiURL, membership.ID, membership.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) RemoveMemberFromList(membership Membership) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/memberships/%d?revision=%d", apiURL, membership.ID, membership.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) FilesForListID(listID uint) (*[]File, error) {
	var resp *http.Response
	var err error

	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/files?list_id=%d", apiURL, listID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	f, err := c.jsonHelper.Unmarshal(b, &[]File{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return f.(*[]File), nil
}

func (c OauthClient) FilesForTaskID(taskID uint) (*[]File, error) {
	var resp *http.Response
	var err error

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	resp, err = c.httpHelper.Get(fmt.Sprintf("%s/files?task_id=%d", apiURL, taskID))

	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	f, err := c.jsonHelper.Unmarshal(b, &[]File{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return f.(*[]File), nil
}

func (c OauthClient) File(fileID uint) (*File, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/files/%d", apiURL, fileID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	f, err := c.jsonHelper.Unmarshal(b, &File{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return f.(*File), nil
}

func (c OauthClient) FilePreview(fileID uint) (*FilePreview, error) {
	resp, err := c.httpHelper.Get(fmt.Sprintf("%s/previews?file_id=%d", apiURL, fileID))
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	f, err := c.jsonHelper.Unmarshal(b, &FilePreview{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return f.(*FilePreview), nil
}

func (c OauthClient) CreateUpload(
	contentType string,
	fileName string,
	fileSize uint,
	partNumber uint,
	md5Sum string,
) (*Upload, error) {

	var body []byte

	if contentType == "" {
		return nil, errors.New("contentType must not be empty")
	}

	if fileName == "" {
		return nil, errors.New("fileName must not be empty")
	}

	if fileSize == 0 {
		return nil, errors.New("fileSize must be non-zero")
	}

	url := fmt.Sprintf("%s/uploads?content_type=%s&file_name=%s&file_size=%d", apiURL, contentType, fileName, fileSize)

	if partNumber != 0 {
		url = fmt.Sprintf("%s&part_number=%d", url, partNumber)
	}

	if md5Sum != "" {
		url = fmt.Sprintf("%s&md5_sum=%s", url, md5Sum)
	}

	resp, err := c.httpHelper.Post(url, body)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &Upload{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return u.(*Upload), nil
}

func (c OauthClient) MarkUploadComplete(uploadID uint) (*Upload, error) {
	resp, err := c.httpHelper.Patch(fmt.Sprintf("%s/uploads/%d?state=finished", apiURL, uploadID), nil)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &Upload{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return u.(*Upload), nil
}

func (c OauthClient) ChunkedUploadPart(uploadID uint, partNumber uint, md5Sum string) (*Upload, error) {

	if uploadID == 0 {
		return nil, errors.New("uploadID must be non-zero")
	}

	if partNumber == 0 {
		return nil, errors.New("partNumber must be non-zero")
	}

	url := fmt.Sprintf("%s/uploads/%d/parts?part_number=%d", apiURL, uploadID, partNumber)

	if md5Sum != "" {
		url = fmt.Sprintf("%s&md5_sum=%s", url, md5Sum)
	}

	resp, err := c.httpHelper.Get(url)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	u, err := c.jsonHelper.Unmarshal(b, &Upload{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return u.(*Upload), nil
}

func (c OauthClient) DestroyFile(file File) error {
	resp, err := c.httpHelper.Delete(fmt.Sprintf("%s/files/%d?revision=%d", apiURL, file.ID, file.Revision))

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	_, err = c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return err
	}
	return nil
}

func (c OauthClient) CreateFile(uploadID uint, taskID uint, localCreatedAt string) (*File, error) {

	if uploadID == 0 {
		return nil, errors.New("uploadID must be > 0")
	}

	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	url := fmt.Sprintf("%s/files?upload_id=%d&task_id=%d", apiURL, uploadID, taskID)

	if localCreatedAt != "" {
		url = fmt.Sprintf("%s&local_created_at=%s", url, localCreatedAt)
	}

	resp, err := c.httpHelper.Post(url, nil)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	b, err := c.readResponseBody(resp)
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}

	f, err := c.jsonHelper.Unmarshal(b, &File{})
	if err != nil {
		c.logger.LogLine(fmt.Sprintf("response: %v", resp))
		return nil, err
	}
	return f.(*File), nil
}
