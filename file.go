package wundergo

type File struct {
	ID             uint   `json:"id"`
	URL            string `json:"url"`
	TaskID         uint   `json:"task_id"`
	ListID         uint   `json:"list_id"`
	UserID         uint   `json:"user_id"`
	FileName       string `json:"file_name"`
	ContentType    string `json:"content_type"`
	FileSize       string `json:"file_size"`
	LocalCreatedAt string `json:"local_created_at"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Type           string `json:"type"`
	Revision       uint   `json:"revision"`
}
