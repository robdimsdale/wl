package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wl"
)

// Folders returns Folders created by the current user.
func (c oauthClient) Folders() ([]wl.Folder, error) {
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

	folders := []wl.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folders)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

type folderCreateConfig struct {
	Title   string `json:"title"`
	ListIDs []uint `json:"list_ids"`
}

// CreateFolder creates a new folder with the provided parameters.
// title and listIDs must both be non-empty
func (c oauthClient) CreateFolder(
	title string,
	listIDs []uint,
) (wl.Folder, error) {
	if title == "" {
		return wl.Folder{}, errors.New("title must be non-empty")
	}

	if listIDs == nil || len(listIDs) == 0 {
		return wl.Folder{}, errors.New("listIDs must be non-nil and non-empty")
	}

	fcc := folderCreateConfig{
		Title:   title,
		ListIDs: listIDs,
	}

	body, err := json.Marshal(fcc)
	if err != nil {
		return wl.Folder{}, err
	}

	reqURL := fmt.Sprintf("%s/folders", c.apiURL)

	req, err := c.newPostRequest(reqURL, body)
	if err != nil {
		return wl.Folder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Folder{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return wl.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	folder := wl.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folder)
	if err != nil {
		return wl.Folder{}, err
	}
	return folder, nil
}

// Folder returns the Folder for the corresponding folderID.
func (c oauthClient) Folder(folderID uint) (wl.Folder, error) {
	if folderID == 0 {
		return wl.Folder{}, errors.New("folderID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/folders/%d",
		c.apiURL,
		folderID,
	)

	req, err := c.newGetRequest(url)
	if err != nil {
		return wl.Folder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Folder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	folder := wl.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&folder)
	if err != nil {
		return wl.Folder{}, err
	}
	return folder, nil
}

// UpdateFolder updates the provided Folder.
func (c oauthClient) UpdateFolder(folder wl.Folder) (wl.Folder, error) {
	body, err := json.Marshal(folder)
	if err != nil {
		return wl.Folder{}, err
	}

	url := fmt.Sprintf(
		"%s/folders/%d",
		c.apiURL,
		folder.ID,
	)

	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wl.Folder{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Folder{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Folder{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedFolder := wl.Folder{}
	err = json.NewDecoder(resp.Body).Decode(&returnedFolder)
	if err != nil {
		return wl.Folder{}, err
	}
	return returnedFolder, nil
}

// DeleteFolder deletes the provided folder.
func (c oauthClient) DeleteFolder(folder wl.Folder) error {
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
func (c oauthClient) FolderRevisions() ([]wl.FolderRevision, error) {
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

	folders := []wl.FolderRevision{}
	err = json.NewDecoder(resp.Body).Decode(&folders)
	if err != nil {
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
	c.logger.Debug(
		"delete-all-folders",
		map[string]interface{}{"folderCount": folderCount},
	)

	idErrChan := make(chan idErr, folderCount)
	for _, f := range folders {
		go func(folder wl.Folder) {
			c.logger.Debug(
				"delete-all-folders - deleting folder",
				map[string]interface{}{"folderID": folder.ID},
			)
			err := c.DeleteFolder(folder)
			idErrChan <- idErr{idType: "folder", id: folder.ID, err: err}
		}(f)
	}

	e := multiIDErr{}
	for i := 0; i < folderCount; i++ {
		idErr := <-idErrChan
		if idErr.err != nil {
			c.logger.Debug(
				"delete-all-folders - error received",
				map[string]interface{}{"id": idErr.id, "err": err},
			)
			e.addError(idErr)
		}
	}

	if len(e.errors()) > 0 {
		return e
	}

	return nil
}
