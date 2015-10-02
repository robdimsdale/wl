package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wl"
)

// Files gets all files for all lists.
func (c oauthClient) Files() ([]wl.File, error) {
	lists, err := c.Lists()
	if err != nil {
		return nil, err
	}

	listCount := len(lists)
	c.logger.Debug(
		"files",
		map[string]interface{}{"listCount": listCount},
	)

	filesChan := make(chan []wl.File, listCount)
	idErrChan := make(chan idErr, listCount)
	for _, l := range lists {
		go func(list wl.List) {
			c.logger.Debug(
				"files - getting files for list",
				map[string]interface{}{"listID": list.ID},
			)
			files, err := c.FilesForListID(list.ID)
			idErrChan <- idErr{idType: "list", id: list.ID, err: err}
			filesChan <- files
		}(l)
	}

	e := multiIDErr{}
	for i := 0; i < listCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"files - error received getting files for list",
				map[string]interface{}{"listID": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	allFiles := []wl.File{}
	for i := 0; i < listCount; i++ {
		files := <-filesChan
		allFiles = append(allFiles, files...)
	}

	if len(e.errors()) > 0 {
		return allFiles, e
	}

	return allFiles, nil
}

// FilesForListID returns the Files associated with the provided List.
func (c oauthClient) FilesForListID(listID uint) ([]wl.File, error) {
	if listID == 0 {
		return nil, errors.New("listID must be > 0")
	}
	url := fmt.Sprintf(
		"%s/files?list_id=%d",
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

	files := []wl.File{}
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// FilesForTaskID returns the Files associated with the provided Task.
func (c oauthClient) FilesForTaskID(taskID uint) ([]wl.File, error) {
	if taskID == 0 {
		return nil, errors.New("taskID must be > 0")
	}
	url := fmt.Sprintf(
		"%s/files?task_id=%d",
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

	files := []wl.File{}
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// File returns the File for the corresponding taskID.
func (c oauthClient) File(fileID uint) (wl.File, error) {
	if fileID == 0 {
		return wl.File{}, errors.New("fileID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/files/%d",
		c.apiURL,
		fileID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wl.File{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.File{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.File{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	task := wl.File{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return wl.File{}, err
	}
	return task, nil
}

// CreateFile creates a File, associating an upload with a task.
func (c oauthClient) CreateFile(uploadID uint, taskID uint) (wl.File, error) {
	if uploadID == 0 {
		return wl.File{}, errors.New("uploadID must be > 0")
	}

	if taskID == 0 {
		return wl.File{}, errors.New("taskID must be > 0")
	}

	body := []byte(fmt.Sprintf(`{"upload_id":%d,"task_id":%d}`, uploadID, taskID))

	url := fmt.Sprintf("%s/files", c.apiURL)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return wl.File{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.File{}, err
	}
	if err != nil {
		return wl.File{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wl.File{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	file := wl.File{}
	err = json.NewDecoder(resp.Body).Decode(&file)
	if err != nil {
		return wl.File{}, err
	}
	return file, nil
}

// DestroyFile deletes the provided File.
func (c oauthClient) DestroyFile(file wl.File) error {
	url := fmt.Sprintf(
		"%s/files/%d?revision=%d",
		c.apiURL,
		file.ID,
		file.Revision,
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
