package wundergo

import "time"

// UploadPart contains information about a partial upload.
type UploadPart struct {
	URL           string `json:"url"`
	Date          string `json:"date"`
	Authorization string `json:"authorization"`
}

// Upload contains information about an upload.
// See also File.
type Upload struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	State     string     `json:"state"`
	ExpiresAt time.Time  `json:"expires_at"`
	Part      UploadPart `json:"part"`
}
