package wundergo

// TaskComment represents information about a comment on a Task.
type TaskComment struct {
	ID        uint   `json:"id"`
	TaskID    uint   `json:"task_id"`
	Revision  uint   `json:"revision"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}
