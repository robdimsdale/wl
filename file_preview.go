package wundergo

type FilePreview struct {
	URL       string `json:"url"`
	Size      string `json:"size"`
	ExpiresAt string `json:"expires_at"`
}
