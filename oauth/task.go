package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/robdimsdale/wl"
)

// Tasks gets all tasks for all lists.
func (c oauthClient) Tasks() ([]wl.Task, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"tasks",
		map[string]interface{}{"listCount": listCount},
	)

	tasksChan := make(chan []wl.Task, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wl.List) {
			c.logger.Debug(
				"tasks - getting tasks for list",
				map[string]interface{}{"listID": list.ID},
			)
			tasks, err := c.TasksForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			tasksChan <- tasks
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"tasks - error received getting tasks for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	totalTasks := []wl.Task{}
	for i := 0; i < listCount; i++ {
		tasks := <-tasksChan
		totalTasks = append(totalTasks, tasks...)
	}

	if len(e.errors()) > 0 {
		return totalTasks, e
	}

	return totalTasks, nil
}

// CompletedTasks returns all tasks filtered by whether they are completed.
func (c oauthClient) CompletedTasks(completed bool) ([]wl.Task, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"tasks",
		map[string]interface{}{"listCount": listCount},
	)

	tasksChan := make(chan []wl.Task, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wl.List) {
			c.logger.Debug(
				"tasks - getting tasks for list",
				map[string]interface{}{"listID": list.ID},
			)
			tasks, err := c.CompletedTasksForListID(list.ID, completed)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			tasksChan <- tasks
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"tasks - error received getting tasks for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	totalTasks := []wl.Task{}
	for i := 0; i < listCount; i++ {
		tasks := <-tasksChan
		totalTasks = append(totalTasks, tasks...)
	}

	if len(e.errors()) > 0 {
		return totalTasks, e
	}

	return totalTasks, nil
}

// TasksForListID returns Tasks for the provided listID.
func (c oauthClient) TasksForListID(listID uint) ([]wl.Task, error) {
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	tasks := []transportTask{}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	return tasksFromTransport(tasks)
}

// CompletedTasksForListID returns tasks filtered by whether they are completed.
func (c oauthClient) CompletedTasksForListID(listID uint, completed bool) ([]wl.Task, error) {
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

	tasks := []transportTask{}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	return tasksFromTransport(tasks)
}

// Task returns the Task for the corresponding taskID.
func (c oauthClient) Task(taskID uint) (wl.Task, error) {
	url := fmt.Sprintf(
		"%s/tasks/%d",
		c.apiURL,
		taskID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wl.Task{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Task{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Task{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	task := transportTask{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return wl.Task{}, err
	}
	return taskFromTransport(task)
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
	Title           string   `json:"title,omitempty"`
	Revision        uint     `json:"revision"`
	AssigneeID      uint     `json:"assignee_id,omitempty"`
	ListID          uint     `json:"list_id,omitempty"`
	Completed       bool     `json:"completed"`
	RecurrenceType  string   `json:"recurrence_type,omitempty"`
	RecurrenceCount uint     `json:"recurrence_count,omitempty"`
	DueDate         string   `json:"due_date,omitempty"`
	Starred         bool     `json:"starred"`
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
	dueDate time.Time,
	starred bool,
) (wl.Task, error) {

	if listID == 0 {
		return wl.Task{}, errors.New("listID must be > 0")
	}

	err := c.validateRecurrence(recurrenceType, recurrenceCount)
	if err != nil {
		return wl.Task{}, err
	}

	tcc := taskCreateConfig{
		ListID:          listID,
		Title:           title,
		AssigneeID:      assigneeID,
		Completed:       completed,
		RecurrenceType:  recurrenceType,
		RecurrenceCount: recurrenceCount,
		DueDate:         dueDateToString(dueDate),
		Starred:         starred,
	}

	body, err := json.Marshal(tcc)
	if err != nil {
		return wl.Task{}, err
	}

	url := fmt.Sprintf("%s/tasks", c.apiURL)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wl.Task{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Task{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wl.Task{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	task := transportTask{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return wl.Task{}, err
	}
	return taskFromTransport(task)
}

// UpdateTask updates the provided Task.
func (c oauthClient) UpdateTask(task wl.Task) (wl.Task, error) {
	err := c.validateRecurrence(task.RecurrenceType, task.RecurrenceCount)
	if err != nil {
		return wl.Task{}, err
	}

	origTask, err := c.Task(task.ID)
	if err != nil {
		return wl.Task{}, err
	}

	tuc := TaskUpdateConfig{
		Title:     task.Title,
		Revision:  task.Revision,
		Completed: task.Completed,
		Starred:   task.Starred,
		ListID:    task.ListID,

		Remove: []string{},
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
		tuc.DueDate = dueDateToString(origTask.DueDate)
	} else {
		if task.DueDate.IsZero() {
			tuc.Remove = append(tuc.Remove, "due_date")
		} else {
			tuc.DueDate = dueDateToString(task.DueDate)
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

	body, err := json.Marshal(tuc)
	if err != nil {
		return wl.Task{}, err
	}

	url := fmt.Sprintf(
		"%s/tasks/%d",
		c.apiURL,
		task.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wl.Task{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Task{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Task{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	transport := transportTask{}
	err = json.NewDecoder(resp.Body).Decode(&transport)
	if err != nil {
		return wl.Task{}, err
	}
	return taskFromTransport(transport)
}

// DeleteTask deletes the provided Task.
func (c oauthClient) DeleteTask(task wl.Task) error {
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

// DeleteAllTasks gets a list of all tasks via Tasks() and deletes them
// via DeleteTask(task)
func (c oauthClient) DeleteAllTasks() error {
	tasks, err := c.Tasks()
	if err != nil {
		return err
	}

	taskCount := len(tasks)
	c.logger.Debug(
		"delete-all-tasks",
		map[string]interface{}{"taskCount": taskCount},
	)

	idErrChan := make(chan idErr, taskCount)
	for _, f := range tasks {
		go func(task wl.Task) {
			c.logger.Debug(
				"delete-all-tasks - deleting task",
				map[string]interface{}{"taskID": task.ID},
			)
			err := c.DeleteTask(task)
			idErrChan <- idErr{idType: "task", id: task.ID, err: err}
		}(f)
	}

	e := multiIDErr{}
	for i := 0; i < taskCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"delete-all-tasks - error received",
				map[string]interface{}{"id": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	if len(e.errors()) > 0 {
		return e
	}

	return nil
}

type transportTask struct {
	ID              uint      `json:"id" yaml:"id"`
	AssigneeID      uint      `json:"assignee_id" yaml:"assignee_id"`
	AssignerID      uint      `json:"assigner_id" yaml:"assigner_id"`
	CreatedAt       time.Time `json:"created_at" yaml:"created_at"`
	CreatedByID     uint      `json:"created_by_id" yaml:"created_by_id"`
	DueDate         string    `json:"due_date" yaml:"due_date"`
	ListID          uint      `json:"list_id" yaml:"list_id"`
	Revision        uint      `json:"revision" yaml:"revision"`
	Starred         bool      `json:"starred" yaml:"starred"`
	Title           string    `json:"title" yaml:"title"`
	Completed       bool      `json:"completed" yaml:"completed"`
	CompletedAt     time.Time `json:"completed_at" yaml:"completed_at"`
	CompletedByID   uint      `json:"completed_by" yaml:"completed_by"`
	RecurrenceType  string    `json:"recurrence_type" yaml:"recurrence_type"`
	RecurrenceCount uint      `json:"recurrence_count" yaml:"recurrence_count"`
}

func tasksFromTransport(transportTasks []transportTask) ([]wl.Task, error) {
	tasks := make([]wl.Task, len(transportTasks))
	var err error

	for i, t := range transportTasks {
		tasks[i], err = taskFromTransport(t)
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

func taskFromTransport(t transportTask) (wl.Task, error) {
	dueDate, err := parseDueDate(t.DueDate, t.CreatedAt.Location())
	if err != nil {
		return wl.Task{}, err
	}

	return wl.Task{
		ID:              t.ID,
		AssigneeID:      t.AssigneeID,
		AssignerID:      t.AssignerID,
		CreatedAt:       t.CreatedAt,
		CreatedByID:     t.CreatedByID,
		DueDate:         dueDate,
		ListID:          t.ListID,
		Revision:        t.Revision,
		Starred:         t.Starred,
		Title:           t.Title,
		Completed:       t.Completed,
		CompletedAt:     t.CompletedAt,
		CompletedByID:   t.CompletedByID,
		RecurrenceType:  t.RecurrenceType,
		RecurrenceCount: t.RecurrenceCount,
	}, nil
}

func parseDueDate(dueDate string, location *time.Location) (time.Time, error) {
	if dueDate == "" {
		return time.Time{}, nil
	}

	splitDate := strings.Split(dueDate, "-")
	if len(splitDate) < 3 {
		return time.Now(), fmt.Errorf("Failed to parse dueDate into expected YYYY-MM-DD format: %s", dueDate)
	}

	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return time.Now(), err
	}

	monthInt, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return time.Now(), err
	}
	month := time.Month(monthInt)

	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		return time.Now(), err
	}

	hour := 0
	minute := 0
	second := 0
	nano := 0

	return time.Date(year, month, day, hour, minute, second, nano, location), nil
}

func dueDateToString(dueDate time.Time) string {
	if (dueDate == time.Time{}) {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", dueDate.Year(), dueDate.Month(), dueDate.Day())
}
