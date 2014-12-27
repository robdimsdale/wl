package wundergo

type Note struct {
	ID        uint   `json:"id"`
	TaskID    uint   `json:"task_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Revision  uint   `json:"revision"`
}
