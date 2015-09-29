package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// Subtasks gets all tasks for all lists.
func (c oauthClient) Subtasks() ([]wundergo.Subtask, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"tasks",
		map[string]interface{}{"listCount": listCount},
	)

	subtasksChan := make(chan []wundergo.Subtask, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wundergo.List) {
			c.logger.Debug(
				"tasks - getting tasks for list",
				map[string]interface{}{"listID": list.ID},
			)
			subtasks, err := c.SubtasksForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			subtasksChan <- subtasks
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

	if len(e.errors()) > 0 {
		return nil, e
	}

	totalSubtasks := []wundergo.Subtask{}
	for i := 0; i < listCount; i++ {
		subtasks := <-subtasksChan
		totalSubtasks = append(totalSubtasks, subtasks...)
	}

	return totalSubtasks, nil
}

// CompletedSubtasks returns all tasks filtered by whether they are completed.
func (c oauthClient) CompletedSubtasks(completed bool) ([]wundergo.Subtask, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"tasks",
		map[string]interface{}{"listCount": listCount},
	)

	subtasksChan := make(chan []wundergo.Subtask, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wundergo.List) {
			c.logger.Debug(
				"subtasks - getting subtasks for list",
				map[string]interface{}{"listID": list.ID},
			)
			subtasks, err := c.CompletedSubtasksForListID(list.ID, completed)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			subtasksChan <- subtasks
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"subtasks - error received getting subtasks for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	if len(e.errors()) > 0 {
		return nil, e
	}

	totalSubtasks := []wundergo.Subtask{}
	for i := 0; i < listCount; i++ {
		subtasks := <-subtasksChan
		totalSubtasks = append(totalSubtasks, subtasks...)
	}

	return totalSubtasks, nil
}

// SubtasksForListID returns the Subtasks associated with the provided listID.
func (c oauthClient) SubtasksForListID(listID uint) ([]wundergo.Subtask, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtasks?list_id=%d",
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

	subtasks := []wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtasks)
	if err != nil {
		return nil, err
	}
	return subtasks, nil
}

// SubtasksForTaskID returns the Subtasks associated with the provided taskID.
func (c oauthClient) SubtasksForTaskID(taskID uint) ([]wundergo.Subtask, error) {
	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtasks?task_id=%d",
		c.apiURL,
		taskID,
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

	subtasks := []wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtasks)
	if err != nil {
		return nil, err
	}
	return subtasks, nil
}

// CompletedSubtasksForListID returns subtasks for the provided List,
// filtered on whether they are completed.
func (c oauthClient) CompletedSubtasksForListID(listID uint, completed bool) ([]wundergo.Subtask, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtasks?list_id=%d&completed=%t",
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

	subtasks := []wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtasks)
	if err != nil {
		return nil, err
	}
	return subtasks, nil
}

// CompletedSubtasksForTaskID returns subtasks for the provided List,
// filtered on whether they are completed.
func (c oauthClient) CompletedSubtasksForTaskID(taskID uint, completed bool) ([]wundergo.Subtask, error) {
	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtasks?task_id=%d&completed=%t",
		c.apiURL,
		taskID,
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

	subtasks := []wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtasks)
	if err != nil {
		return nil, err
	}
	return subtasks, nil
}

// Subtask returns the subtask for the corresponding subtaskID.
func (c oauthClient) Subtask(subtaskID uint) (wundergo.Subtask, error) {
	if subtaskID == 0 {
		return wundergo.Subtask{}, errors.New("subtaskID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtasks/%d",
		c.apiURL,
		subtaskID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Subtask{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Subtask{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Subtask{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtask := wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtask)
	if err != nil {
		return wundergo.Subtask{}, err
	}
	return subtask, nil
}

// CreateSubtask creates a Subtask for the provided parameters.
func (c oauthClient) CreateSubtask(
	title string,
	taskID uint,
	completed bool,
) (wundergo.Subtask, error) {

	if taskID == 0 {
		return wundergo.Subtask{}, errors.New("taskID must be > 0")
	}

	body := []byte(fmt.Sprintf(
		`{"title":"%s","task_id":%d,"completed":%t}`,
		title,
		taskID,
		completed,
	))

	url := fmt.Sprintf("%s/subtasks", c.apiURL)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wundergo.Subtask{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Subtask{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wundergo.Subtask{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	subtask := wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtask)
	if err != nil {
		return wundergo.Subtask{}, err
	}
	return subtask, nil
}

// UpdateSubtask updates the provided Subtask.
func (c oauthClient) UpdateSubtask(subtask wundergo.Subtask) (wundergo.Subtask, error) {
	body, err := json.Marshal(subtask)
	if err != nil {
		return wundergo.Subtask{}, err
	}

	url := fmt.Sprintf(
		"%s/subtasks/%d",
		c.apiURL,
		subtask.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wundergo.Subtask{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Subtask{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Subtask{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedSubtask := wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&returnedSubtask)
	if err != nil {
		return wundergo.Subtask{}, err
	}
	return returnedSubtask, nil
}

// DeleteSubtask deletes the provided Subtask.
func (c oauthClient) DeleteSubtask(subtask wundergo.Subtask) error {
	url := fmt.Sprintf(
		"%s/subtasks/%d?revision=%d",
		c.apiURL,
		subtask.ID,
		subtask.Revision,
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
