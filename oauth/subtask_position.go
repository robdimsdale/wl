package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pivotal-golang/lager"
	"github.com/robdimsdale/wundergo"
)

// SubtaskPositionsForListID returns the positions of all Subtasks in the List
// associated with the provided listID.
// The returned SubtaskPosition.Values might be empty if the Subtasks have never been reordered.
func (c oauthClient) SubtaskPositionsForListID(listID uint) ([]wundergo.Position, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtask_positions?list_id=%d",
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

	subtaskPositions := []wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(&subtaskPositions)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return nil, err
	}
	return subtaskPositions, nil
}

// SubtaskPositionsForTaskID returns the positions of all Subtasks in the Task
// associated with the provided taskID.
// The returned SubtaskPosition.Values might be empty if the Subtasks have never been reordered.
func (c oauthClient) SubtaskPositionsForTaskID(taskID uint) ([]wundergo.Position, error) {
	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtask_positions?task_id=%d",
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

	subtaskPositions := []wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(&subtaskPositions)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return nil, err
	}
	return subtaskPositions, nil
}

// SubtaskPosition returns the SubtaskPosition associated with the provided subtaskPositionID.
func (c oauthClient) SubtaskPosition(subTaskPositionID uint) (wundergo.Position, error) {
	if subTaskPositionID == 0 {
		return wundergo.Position{}, errors.New("subTaskPositionID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtask_positions/%d",
		c.apiURL,
		subTaskPositionID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Position{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtaskPosition := wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(&subtaskPosition)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Position{}, err
	}
	return subtaskPosition, nil
}

// UpdateSubtaskPosition updates the provided SubtaskPosition.
// This will reorder the Subtasks.
func (c oauthClient) UpdateSubtaskPosition(subTaskPosition wundergo.Position) (wundergo.Position, error) {
	body, err := json.Marshal(subTaskPosition)
	if err != nil {
		return wundergo.Position{}, err
	}

	url := fmt.Sprintf(
		"%s/subtask_positions/%d",
		c.apiURL,
		subTaskPosition.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wundergo.Position{}, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return wundergo.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedSubtaskPosition := wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(&returnedSubtaskPosition)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Position{}, err
	}
	return returnedSubtaskPosition, nil
}
