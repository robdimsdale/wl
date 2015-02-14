package wundergo

import "time"

// File contains the information about an uploaded file.
// See also Upload.
type File struct {
	ID             uint      `json:"id"`
	URL            string    `json:"url"`
	TaskID         uint      `json:"task_id"`
	ListID         uint      `json:"list_id"`
	UserID         uint      `json:"user_id"`
	FileName       string    `json:"file_name"`
	ContentType    string    `json:"content_type"`
	FileSize       string    `json:"file_size"`
	LocalCreatedAt time.Time `json:"local_created_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Type           string    `json:"type"`
	Revision       uint      `json:"revision"`
}
