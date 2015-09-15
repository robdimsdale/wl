package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

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
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
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

	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtasks := []wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtasks)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
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

	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtasks := []wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtasks)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
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

	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtasks := []wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtasks)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
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
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
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
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
		return wundergo.Subtask{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wundergo.Subtask{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	subtask := wundergo.Subtask{}
	err = json.NewDecoder(resp.Body).Decode(&subtask)
	if err != nil {
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
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
		c.logger.Debug("", map[string]interface{}{"response": newLoggableResponse(resp)})
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
