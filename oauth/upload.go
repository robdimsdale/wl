package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/robdimsdale/wl"
)

// UploadFile uploads the local file to Wunderlist.
// Currently it does not support multi-part uploading: the entire file
// must fit into a single part
// md5sum is optional; other fields are required
func (c oauthClient) UploadFile(
	localFilePath string,
	remoteFileName string,
	contentType string,
	md5sum string,
) (wl.Upload, error) {
	fileContents, err := c.readLocalFile(localFilePath)
	if err != nil {
		return wl.Upload{}, err
	}

	initialUploadResp, err := c.createUpload(remoteFileName, contentType, len(fileContents), md5sum)
	uploadID := initialUploadResp.ID

	// Upload actual data using returned URL
	// Do not worry about multi-part upload for now
	err = c.uploadAPart(initialUploadResp.Part, fileContents)
	if err != nil {
		return wl.Upload{}, err
	}

	return c.finishUpload(uploadID)
}

func (c oauthClient) createUpload(
	remoteFileName string,
	contentType string,
	fileSize int,
	md5sum string,
) (uploadResponse, error) {
	if remoteFileName == "" {
		return uploadResponse{}, errors.New("remoteFileName must be non-empty")
	}

	url := fmt.Sprintf("%s/uploads", c.apiURL)

	bodyString := fmt.Sprintf(
		`{"content_type":"%s","file_name":"%s","file_size":%d`,
		contentType,
		remoteFileName,
		fileSize,
	)

	if md5sum != "" {
		bodyString = fmt.Sprintf(`%s,"md5sum":"%s"`, bodyString, md5sum)
	}

	bodyString = fmt.Sprintf("%s}", bodyString)

	body := []byte(bodyString)

	req, err := c.newPostRequest(url, body)
	if err != nil {
		return uploadResponse{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return uploadResponse{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return uploadResponse{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusCreated)
	}

	uploadResp := uploadResponse{}
	err = json.NewDecoder(resp.Body).Decode(&uploadResp)
	if err != nil {
		return uploadResponse{}, err
	}

	return uploadResp, err
}

type uploadResponse struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	State     string     `json:"state"`
	Part      uploadPart `json:"part"`
	ExpiresAt string     `json:"expires_at"`
}

type uploadPart struct {
	URL           string `json:"url"`
	Date          string `json:"date"`
	Authorization string `json:"authorization"`
}

func (c oauthClient) readLocalFile(localFilePath string) ([]byte, error) {
	// Upload actual data using returned URL
	// Do not worry about multi-part upload for now
	c.logger.Debug("reading local file", map[string]interface{}{"localFilePath": localFilePath})
	fileContents, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		return nil, err
	}

	return fileContents, nil
}

func (c oauthClient) uploadAPart(part uploadPart, fileContents []byte) error {
	req, err := http.NewRequest("PUT", part.URL, nil)
	if err != nil {
		return err
	}

	c.addBody(req, fileContents)
	req.Header.Add("Content-Type", "")
	req.Header.Add("x-amz-date", part.Date)
	req.Header.Add("Authorization", part.Authorization)

	c.logger.Debug(" - posting local file contents", map[string]interface{}{"URL": part.URL})
	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	return nil
}

func (c oauthClient) finishUpload(uploadID uint) (wl.Upload, error) {
	// Mark upload as finished
	c.logger.Debug(" - marking upload as finished", map[string]interface{}{"uploadID": uploadID})
	url := fmt.Sprintf("%s/uploads/%d", c.apiURL, uploadID)
	body := []byte(`{"state":"finished"}`)
	req, err := c.newPatchRequest(url, body)
	if err != nil {
		return wl.Upload{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wl.Upload{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wl.Upload{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	returnedUpload := wl.Upload{}
	err = json.NewDecoder(resp.Body).Decode(&returnedUpload)
	if err != nil {
		return wl.Upload{}, err
	}

	return returnedUpload, nil
}
