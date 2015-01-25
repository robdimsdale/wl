package wundergo

type UploadPart struct {
	URL           string `json:"url"`
	Date          string `json:"date"`
	Authorization string `json:"authorization"`
}

type Upload struct {
	ID        uint       `json:"id"`
	UserID    uint       `json:"user_id"`
	State     string     `json:"state"`
	ExpiresAt string     `json:"expires_at"`
	Part      UploadPart `json:"part"`
}
