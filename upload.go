package wundergo

// Upload contains information about uploads.
// Uploads represent uploaded files.
type Upload struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	State  string `json:"state"`
}
