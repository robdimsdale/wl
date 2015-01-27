package wundergo

// FilePreview contains information about an image thumnail for a file.
type FilePreview struct {
	URL       string `json:"url"`
	Size      string `json:"size"`
	ExpiresAt string `json:"expires_at"`
}
