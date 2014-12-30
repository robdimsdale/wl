package wundergo

type Reminder struct {
	ID        uint   `json:"id"`
	Date      string `json:"date"`
	TaskID    uint   `json:"task_id"`
	Revision  uint   `json:"revision"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
