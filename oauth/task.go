package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// TasksForListID returns Tasks for the provided listID.
func (c oauthClient) TasksForListID(listID uint) ([]wundergo.Task, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}
	url := fmt.Sprintf(
		"%s/tasks?list_id=%d",
		c.apiURL,
		listID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	tasks := []wundergo.Task{}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return nil, err
	}
	return tasks, nil
}

// CompletedTasksForListID returns tasks filtered by whether they are completed.
func (c oauthClient) CompletedTasksForListID(listID uint, completed bool) ([]wundergo.Task, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/tasks?list_id=%d&completed=%t",
		c.apiURL,
		listID,
		completed,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	tasks := []wundergo.Task{}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return nil, err
	}
	return tasks, nil
}

// Task returns the Task for the corresponding taskID.
func (c oauthClient) Task(taskID uint) (wundergo.Task, error) {
	url := fmt.Sprintf(
		"%s/tasks/%d",
		c.apiURL,
		taskID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Task{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Task{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Task{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	task := wundergo.Task{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return wundergo.Task{}, err
	}
	return task, nil
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

// TaskUpdateConfig contains information required to update an existing task.
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

// CreateTask creates a task with the provided parameters.
func (c oauthClient) CreateTask(
	title string,
	listID uint,
	assigneeID uint,
	completed bool,
	recurrenceType string,
	recurrenceCount uint,
	dueDate string,
	starred bool,
) (wundergo.Task, error) {

	if listID == 0 {
		return wundergo.Task{}, errors.New("listID must be > 0")
	}

	err := c.validateRecurrence(recurrenceType, recurrenceCount)
	if err != nil {
		return wundergo.Task{}, err
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

	body, err := json.Marshal(tcc)
	if err != nil {
		return wundergo.Task{}, err
	}

	url := fmt.Sprintf("%s/tasks", c.apiURL)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wundergo.Task{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Task{}, err
	}
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return wundergo.Task{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wundergo.Task{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	task := wundergo.Task{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return wundergo.Task{}, err
	}
	return task, nil
}

// UpdateTask updates the provided Task.
func (c oauthClient) UpdateTask(task wundergo.Task) (wundergo.Task, error) {
	err := c.validateRecurrence(task.RecurrenceType, task.RecurrenceCount)
	if err != nil {
		return wundergo.Task{}, err
	}

	origTask, err := c.Task(task.ID)
	if err != nil {
		return wundergo.Task{}, err
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

	if origTask.RecurrenceCount == task.RecurrenceCount &&
		origTask.RecurrenceType == task.RecurrenceType {
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

	body, err := json.Marshal(tuc)
	if err != nil {
		return wundergo.Task{}, err
	}

	url := fmt.Sprintf(
		"%s/tasks/%d",
		c.apiURL,
		task.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wundergo.Task{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Task{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Task{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedTask := wundergo.Task{}
	err = json.NewDecoder(resp.Body).Decode(&returnedTask)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return wundergo.Task{}, err
	}
	return returnedTask, nil
}

// DeleteTask deletes the provided Task.
func (c oauthClient) DeleteTask(task wundergo.Task) error {
	url := fmt.Sprintf(
		"%s/tasks/%d?revision=%d",
		c.apiURL,
		task.ID,
		task.Revision,
	)

	req, err := c.newDeleteRequest(url)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	return nil
}
