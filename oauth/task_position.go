package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// TaskPositions gets all tasks for all lists.
func (c oauthClient) TaskPositions() ([]wundergo.Position, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"taskPositions",
		map[string]interface{}{"listCount": listCount},
	)

	taskPositionsChan := make(chan []wundergo.Position, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wundergo.List) {
			c.logger.Debug(
				"taskPositions - getting taskPositions for list",
				map[string]interface{}{"listID": list.ID},
			)
			taskPositions, err := c.TaskPositionsForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			taskPositionsChan <- taskPositions
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"taskPositions - error received getting taskPositions for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	totalTaskPositions := []wundergo.Position{}
	for i := 0; i < listCount; i++ {
		taskPositions := <-taskPositionsChan
		totalTaskPositions = append(totalTaskPositions, taskPositions...)
	}

	if len(e.errors()) > 0 {
		return totalTaskPositions, e
	}

	return totalTaskPositions, nil
}

// TaskPositionsForListID returns the positions of all Tasks in the List
// associated with the provided listID.
// The returned TaskPosition.Values might be empty if the Tasks have never been reordered.
func (c oauthClient) TaskPositionsForListID(listID uint) ([]wundergo.Position, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/task_positions?list_id=%d",
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

	taskPositions := []wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(&taskPositions)
	if err != nil {
		return nil, err
	}
	return taskPositions, nil
}

// TaskPosition returns the TaskPosition associated with the provided taskPositionID.
func (c oauthClient) TaskPosition(taskPositionID uint) (wundergo.Position, error) {
	if taskPositionID == 0 {
		return wundergo.Position{}, errors.New("taskPositionID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/task_positions/%d",
		c.apiURL,
		taskPositionID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Position{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	taskPosition := wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(&taskPosition)
	if err != nil {
		return wundergo.Position{}, err
	}
	return taskPosition, nil
}

// UpdateTaskPosition updates the provided TaskPosition.
// This will reorder the Tasks.
func (c oauthClient) UpdateTaskPosition(taskPosition wundergo.Position) (wundergo.Position, error) {
	body, err := json.Marshal(taskPosition)
	if err != nil {
		return wundergo.Position{}, err
	}

	url := fmt.Sprintf(
		"%s/task_positions/%d",
		c.apiURL,
		taskPosition.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wundergo.Position{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Position{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Position{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedTaskPosition := wundergo.Position{}
	err = json.NewDecoder(resp.Body).Decode(&returnedTaskPosition)
	if err != nil {
		return wundergo.Position{}, err
	}
	return returnedTaskPosition, nil
}
