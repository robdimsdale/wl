package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/robdimsdale/wundergo"
)

// FilePreview returns the FilePreview for the corresponding fileID.
// fileID must be > 0; platform and size are included if they are non=empty,
// and are not valididated otherwise.
func (c oauthClient) FilePreview(
	fileID uint,
	platform string,
	size string,
) (wundergo.FilePreview, error) {
	if fileID == 0 {
		return wundergo.FilePreview{}, errors.New("fileID must be > 0")
	}

	url := fmt.Sprintf(
		"%s/previews?file_id=%d",
		c.apiURL,
		fileID,
	)

	if platform != "" {
		url = fmt.Sprintf(
			"%s&platform=%s",
			url,
			platform,
		)
	}

	if size != "" {
		url = fmt.Sprintf(
			"%s&size=%s",
			url,
			size,
		)
	}

	req, err := c.newGetRequest(url)
	if err != nil {
		return wundergo.FilePreview{}, err
	}

	resp, err := c.do(req)
	if err != nil {
		return wundergo.FilePreview{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return wundergo.FilePreview{}, fmt.Errorf("Unexpected response code %d - expected %d", resp.StatusCode, http.StatusOK)
	}

	task := wundergo.FilePreview{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return wundergo.FilePreview{}, err
	}
	return task, nil
}
