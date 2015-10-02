package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wl"
)

// SubtaskPositions gets all subtask positions for all lists.
func (c oauthClient) SubtaskPositions() ([]wl.Position, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"subtaskPositions",
		map[string]interface{}{"listCount": listCount},
	)

	subtaskPositionsChan := make(chan []wl.Position, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wl.List) {
			c.logger.Debug(
				"subtaskPositions - getting subtaskPositions for list",
				map[string]interface{}{"listID": list.ID},
			)
			subtaskPositions, err := c.SubtaskPositionsForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			subtaskPositionsChan <- subtaskPositions
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"subtaskPositions - error received getting subtaskPositions for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	totalSubtaskPositions := []wl.Position{}
	for i := 0; i < listCount; i++ {
		subtaskPositions := <-subtaskPositionsChan
		totalSubtaskPositions = append(totalSubtaskPositions, subtaskPositions...)
	}

	if len(e.errors()) > 0 {
		return totalSubtaskPositions, e
	}

	return totalSubtaskPositions, nil
}

// SubtaskPositionsForListID returns the positions of all Subtasks in the List
// associated with the provided listID.
// The returned SubtaskPosition.Values might be empty if the Subtasks have never been reordered.
func (c oauthClient) SubtaskPositionsForListID(listID uint) ([]wl.Position, error) {
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

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtaskPositions := []wl.Position{}
	err = json.NewDecoder(resp.Body).Decode(&subtaskPositions)
	if err != nil {
		return nil, err
	}
	return subtaskPositions, nil
}

// SubtaskPositionsForTaskID returns the positions of all Subtasks in the Task
// associated with the provided taskID.
// The returned SubtaskPosition.Values might be empty if the Subtasks have never been reordered.
func (c oauthClient) SubtaskPositionsForTaskID(taskID uint) ([]wl.Position, error) {
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

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtaskPositions := []wl.Position{}
	err = json.NewDecoder(resp.Body).Decode(&subtaskPositions)
	if err != nil {
		return nil, err
	}
	return subtaskPositions, nil
}

// SubtaskPosition returns the SubtaskPosition associated with the provided subtaskPositionID.
func (c oauthClient) SubtaskPosition(subTaskPositionID uint) (wl.Position, error) {
	if subTaskPositionID == 0 {
		return wl.Position{}, errors.New("subTaskPositionID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/subtask_positions/%d",
		c.apiURL,
		subTaskPositionID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wl.Position{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	subtaskPosition := wl.Position{}
	err = json.NewDecoder(resp.Body).Decode(&subtaskPosition)
	if err != nil {
		return wl.Position{}, err
	}
	return subtaskPosition, nil
}

// UpdateSubtaskPosition updates the provided SubtaskPosition.
// This will reorder the Subtasks.
func (c oauthClient) UpdateSubtaskPosition(subTaskPosition wl.Position) (wl.Position, error) {
	body, err := json.Marshal(subTaskPosition)
	if err != nil {
		return wl.Position{}, err
	}

	url := fmt.Sprintf(
		"%s/subtask_positions/%d",
		c.apiURL,
		subTaskPosition.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wl.Position{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedSubtaskPosition := wl.Position{}
	err = json.NewDecoder(resp.Body).Decode(&returnedSubtaskPosition)
	if err != nil {
		return wl.Position{}, err
	}
	return returnedSubtaskPosition, nil
}
