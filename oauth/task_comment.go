package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pivotal-golang/lager"
	"github.com/robdimsdale/wundergo"
)

// TaskCommentsForListID returns TaskComments for the provided listID.
func (c oauthClient) TaskCommentsForListID(listID uint) ([]wundergo.TaskComment, error) {
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

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	taskComments := []wundergo.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComments)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": resp})
		return nil, err
	}
	return taskComments, nil
}

// TaskCommentsForTaskID returns TaskComments for the provided taskID.
func (c oauthClient) TaskCommentsForTaskID(taskID uint) ([]wundergo.TaskComment, error) {
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

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	taskComments := []wundergo.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComments)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": resp})
		return nil, err
	}
	return taskComments, nil
}

// CreateTaskComment creates a TaskComment with the provided content associated with the
// Task for the corresponding taskID.
func (c oauthClient) CreateTaskComment(text string, taskID uint) (wundergo.TaskComment, error) {
	if taskID == 0 {
		return wundergo.TaskComment{}, errors.New("taskID must be > 0")
	}

	body := []byte(fmt.Sprintf(`{"text":"%s","task_id":%d}`, text, taskID))

	url := fmt.Sprintf("%s/task_comments", c.apiURL)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wundergo.TaskComment{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.TaskComment{}, err
	}
	if err != nil {
		c.logger.Debug("", lager.Data{"response": resp})
		return wundergo.TaskComment{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wundergo.TaskComment{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	taskComment := wundergo.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComment)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": resp})
		return wundergo.TaskComment{}, err
	}
	return taskComment, nil
}

// TaskComment returns the TaskComment for the corresponding taskCommentID.
func (c oauthClient) TaskComment(taskCommentID uint) (wundergo.TaskComment, error) {
	if taskCommentID == 0 {
		return wundergo.TaskComment{}, errors.New("taskCommentID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/task_comments/%d",
		c.apiURL,
		taskCommentID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.TaskComment{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.TaskComment{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.TaskComment{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	taskComment := wundergo.TaskComment{}
	err = json.NewDecoder(resp.Body).Decode(&taskComment)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": resp})
		return wundergo.TaskComment{}, err
	}
	return taskComment, nil
}

// DeleteTaskComment deletes the provided TaskComment.
func (c oauthClient) DeleteTaskComment(taskComment wundergo.TaskComment) error {
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

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusNoContent)
	}

	return nil
}
