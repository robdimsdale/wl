package wundergo

// Subtask contains information about a subtask.
// Subtasks are children of tasks.
type Subtask struct {
	ID            uint   `json:"id"`
	TaskID        uint   `json:"task_id"`
	CreatedAt     string `json:"created_at"`
	CreatedByID   uint   `json:"created_by_id"`
	Revision      uint   `json:"revision"`
	Title         string `json:"title"`
	CompletedAt   string `json:"completed_at"`
	CompletedByID uint   `json:"completed_by"`
}
