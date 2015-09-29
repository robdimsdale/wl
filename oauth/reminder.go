package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// Reminders gets all reminders for all lists.
func (c oauthClient) Reminders() ([]wundergo.Reminder, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"reminders",
		map[string]interface{}{"listCount": listCount},
	)

	remindersChan := make(chan []wundergo.Reminder, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wundergo.List) {
			c.logger.Debug(
				"reminders - getting reminders for list",
				map[string]interface{}{"listID": list.ID},
			)
			reminders, err := c.RemindersForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			remindersChan <- reminders
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"reminders - error received getting reminders for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	totalReminders := []wundergo.Reminder{}
	for i := 0; i < listCount; i++ {
		reminders := <-remindersChan
		totalReminders = append(totalReminders, reminders...)
	}

	if len(e.errors()) > 0 {
		return totalReminders, e
	}

	return totalReminders, nil
}

// RemindersForListID returns the Reminders for the List associated with the
// provided listID.
func (c oauthClient) RemindersForListID(listID uint) ([]wundergo.Reminder, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/reminders?list_id=%d",
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

	reminders := []wundergo.Reminder{}
	err = json.NewDecoder(resp.Body).Decode(&reminders)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

// RemindersForTaskID returns the Reminders for the Task associated with the
// provided taskID.
func (c oauthClient) RemindersForTaskID(taskID uint) ([]wundergo.Reminder, error) {
	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/reminders?task_id=%d",
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

	reminders := []wundergo.Reminder{}
	err = json.NewDecoder(resp.Body).Decode(&reminders)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

// Reminder returns the Reminder associated with the provided reminderID.
func (c oauthClient) Reminder(reminderID uint) (wundergo.Reminder, error) {
	url := fmt.Sprintf(
		"%s/reminders/%d",
		c.apiURL,
		reminderID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Reminder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Reminder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Reminder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	reminder := wundergo.Reminder{}
	err = json.NewDecoder(resp.Body).Decode(&reminder)
	if err != nil {
		return wundergo.Reminder{}, err
	}
	return reminder, nil
}

// CreateReminder creates a Reminder with the provided parameters.
func (c oauthClient) CreateReminder(
	date string,
	taskID uint,
	createdByDeviceUdid string,
) (wundergo.Reminder, error) {

	if taskID == 0 {
		return wundergo.Reminder{}, errors.New("taskID must be > 0")
	}

	var body []byte
	if createdByDeviceUdid == "" {
		body = []byte(fmt.Sprintf(`{"date":"%s","task_id":%d}`, date, taskID))
	} else {
		body = []byte(fmt.Sprintf(`{"date":"%s","task_id":%d,"created_by_device_udid":"%s"}`, date, taskID, createdByDeviceUdid))
	}

	url := fmt.Sprintf("%s/reminders", c.apiURL)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wundergo.Reminder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Reminder{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wundergo.Reminder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	reminder := wundergo.Reminder{}
	err = json.NewDecoder(resp.Body).Decode(&reminder)
	if err != nil {
		return wundergo.Reminder{}, err
	}
	return reminder, nil
}

// UpdateReminder updates the provided Reminder.
func (c oauthClient) UpdateReminder(reminder wundergo.Reminder) (wundergo.Reminder, error) {
	body, err := json.Marshal(reminder)
	if err != nil {
		return wundergo.Reminder{}, err
	}

	url := fmt.Sprintf(
		"%s/reminders/%d",
		c.apiURL,
		reminder.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wundergo.Reminder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Reminder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Reminder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedReminder := wundergo.Reminder{}
	err = json.NewDecoder(resp.Body).Decode(&returnedReminder)
	if err != nil {
		return wundergo.Reminder{}, err
	}
	return returnedReminder, nil
}

// DeleteReminder deletes the provided Reminder.
func (c oauthClient) DeleteReminder(reminder wundergo.Reminder) error {
	url := fmt.Sprintf(
		"%s/reminders/%d?revision=%d",
		c.apiURL,
		reminder.ID,
		reminder.Revision,
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
