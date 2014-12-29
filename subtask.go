package wundergo

type Subtask struct {
	ID          uint   `json:"id"`
	TaskID      uint   `json:"task_id"`
	CreatedAt   string `json:"created_at"`
	CreatedByID uint   `json:"created_by_id"`
	Revision    uint   `json:"revision"`
	Title       string `json:"title"`
}
