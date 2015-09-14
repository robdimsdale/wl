package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-golang/lager"
	"github.com/robdimsdale/wundergo"
)

// Folders returns Folders created by the current user.
func (c oauthClient) Folders() ([]wundergo.Folder, error) {
	url := fmt.Sprintf(
		"%s/folders",
		c.apiURL,
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

	folders := []wundergo.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folders)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return nil, err
	}
	return folders, nil
}

type folderCreateConfig struct {
	Title   string `json:"title"`
	ListIDs []uint `json:"list_ids"`
}

// CreateFolder creates a new folder with the provided parameters.
func (c oauthClient) CreateFolder(
	title string,
	listIDs []uint,
) (wundergo.Folder, error) {
	if title == "" {
		return wundergo.Folder{}, errors.New("title must be non-empty")
	}

	if listIDs == nil {
		return wundergo.Folder{}, errors.New("listIDs must be non-nil")
	}

	fcc := folderCreateConfig{
		Title:   title,
		ListIDs: listIDs,
	}

	body, err := json.Marshal(fcc)
	if err != nil {
		return wundergo.Folder{}, err
	}

	reqURL := fmt.Sprintf("%s/folders", c.apiURL)

	req, err := c.newPostRequest(reqURL, body)
	if err != nil {
		return wundergo.Folder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Folder{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.Body != nil {
			b, _ := ioutil.ReadAll(resp.Body)
			c.logger.Debug("", lager.Data{"response.Body": string(b)})
		}
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	folder := wundergo.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folder)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Folder{}, err
	}
	return folder, nil
}

// Folder returns the Folder for the corresponding folderID.
func (c oauthClient) Folder(folderID uint) (wundergo.Folder, error) {
	if folderID == 0 {
		return wundergo.Folder{}, errors.New("folderID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/folders/%d",
		c.apiURL,
		folderID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.Folder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Folder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	folder := wundergo.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folder)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Folder{}, err
	}
	return folder, nil
}

// UpdateFolder updates the provided Folder.
func (c oauthClient) UpdateFolder(folder wundergo.Folder) (wundergo.Folder, error) {
	body, err := json.Marshal(folder)
	if err != nil {
		return wundergo.Folder{}, err
	}

	url := fmt.Sprintf(
		"%s/folders/%d",
		c.apiURL,
		folder.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wundergo.Folder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.Folder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedFolder := wundergo.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&returnedFolder)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return wundergo.Folder{}, err
	}
	return returnedFolder, nil
}

// DeleteFolder deletes the provided folder.
func (c oauthClient) DeleteFolder(folder wundergo.Folder) error {
	url := fmt.Sprintf(
		"%s/folders/%d?revision=%d",
		c.apiURL,
		folder.ID,
		folder.Revision,
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

// FolderRevisions returns FolderRevisions created by the current user.
func (c oauthClient) FolderRevisions() ([]wundergo.FolderRevision, error) {
	url := fmt.Sprintf(
		"%s/folder_revisions",
		c.apiURL,
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

	folders := []wundergo.FolderRevision{}
	err = json.NewDecoder(resp.Body).Decode(&folders)
	if err != nil {
		c.logger.Debug("", lager.Data{"response": newLoggableResponse(resp)})
		return nil, err
	}
	return folders, nil
}

// DeleteAllFolders gets a list of all folders via Folders() and deletes them
// via DeleteFolder(folderID)
func (c oauthClient) DeleteAllFolders() error {
	folders, err := c.Folders()
	if err != nil {
		return err
	}

	folderCount := len(folders)
	c.logger.Debug("delete-all-folders", lager.Data{"folderCount": folderCount})
	idErrChan := make(chan idErr, folderCount)
	for _, f := range folders {
		go func(folder wundergo.Folder) {
			c.logger.Debug("delete-all-folders - deleting folder", lager.Data{"folderID": folder.ID})
			err := c.DeleteFolder(folder)
			idErrChan <- idErr{id: folder.ID, err: err}
		}(f)
	}

	e := multiIDErr{}
	for i := 0; i < len(folders); i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug("delete-all-folders - error received", lager.Data{"id": idErr.id, "err": err})
			e.addError(idErr)
		}
	}

	if len(e.errors()) > 0 {
		return e
	}

	return nil
}
