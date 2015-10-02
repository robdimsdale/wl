package wl

// Upload contains information about uploads.
// Uploads represent uploaded files.
type Upload struct {
	ID     uint   `json:"id" yaml:"id"`
	UserID uint   `json:"user_id" yaml:"user_id"`
	State  string `json:"state" yaml:"state"`
}
