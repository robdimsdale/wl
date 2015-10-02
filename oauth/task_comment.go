package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wl"
)

// TaskComments gets all taskComments for all lists.
func (c oauthClient) TaskComments() ([]wl.TaskComment, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"taskComments",
		map[string]interface{}{"listCount": listCount},
	)

	taskCommentsChan := make(chan []wl.TaskComment, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wl.List) {
			c.logger.Debug(
				"taskComments - getting taskComments for list",
				map[string]interface{}{"listID": list.ID},
			)
			taskComments, err := c.TaskCommentsForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			taskCommentsChan <- taskComments
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"taskComments - error received getting taskComments for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	totalTaskComments := []wl.TaskComment{}
	for i := 0; i < listCount; i++ {
		taskComments := <-taskCommentsChan
		totalTaskComments = append(totalTaskComments, taskComments...)
	}

	if len(e.errors()) > 0 {
		return totalTaskComments, e
	}

	return totalTaskComments, nil
}

// TaskCommentsForListID returns TaskComments for the provided listID.
func (c oauthClient) TaskCommentsForListID(listID uint) ([]wl.TaskComment, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/task_comments?list_id=%d",
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

	taskComments := []wl.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComments)
	if err != nil {
		return nil, err
	}
	return taskComments, nil
}

// TaskCommentsForTaskID returns TaskComments for the provided taskID.
func (c oauthClient) TaskCommentsForTaskID(taskID uint) ([]wl.TaskComment, error) {
	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/task_comments?task_id=%d",
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

	taskComments := []wl.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComments)
	if err != nil {
		return nil, err
	}
	return taskComments, nil
}

// CreateTaskComment creates a TaskComment with the provided content associated with the
// Task for the corresponding taskID.
func (c oauthClient) CreateTaskComment(text string, taskID uint) (wl.TaskComment, error) {
	if taskID == 0 {
		return wl.TaskComment{}, errors.New("taskID must be > 0")
	}

	body := []byte(fmt.Sprintf(`{"text":"%s","task_id":%d}`, text, taskID))

	url := fmt.Sprintf("%s/task_comments", c.apiURL)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wl.TaskComment{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.TaskComment{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wl.TaskComment{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	taskComment := wl.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComment)
	if err != nil {
		return wl.TaskComment{}, err
	}
	return taskComment, nil
}

// TaskComment returns the TaskComment for the corresponding taskCommentID.
func (c oauthClient) TaskComment(taskCommentID uint) (wl.TaskComment, error) {
	if taskCommentID == 0 {
		return wl.TaskComment{}, errors.New("taskCommentID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/task_comments/%d",
		c.apiURL,
		taskCommentID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wl.TaskComment{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.TaskComment{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.TaskComment{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	taskComment := wl.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComment)
	if err != nil {
		return wl.TaskComment{}, err
	}
	return taskComment, nil
}

// DeleteTaskComment deletes the provided TaskComment.
func (c oauthClient) DeleteTaskComment(taskComment wl.TaskComment) error {
	url := fmt.Sprintf(
		"%s/task_comments/%d?revision=%d",
		c.apiURL,
		taskComment.ID,
		taskComment.Revision,
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
