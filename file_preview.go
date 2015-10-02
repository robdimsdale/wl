package wl

import "time"

// FilePreview contains the information about an image thumbnail
type FilePreview struct {
	URL       string    `json:"url" yaml:"url"`
	Size      string    `json:"size" yaml:"size"`
	ExpiresAt time.Time `json:"expires_at" yaml:"expires_at"`
}
