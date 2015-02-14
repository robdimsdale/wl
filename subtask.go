package wundergo

import "time"

// Subtask contains information about a subtask.
// Subtasks are children of tasks.
type Subtask struct {
	ID            uint      `json:"id"`
	TaskID        uint      `json:"task_id"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedByID   uint      `json:"created_by_id"`
	Revision      uint      `json:"revision"`
	Title         string    `json:"title"`
	CompletedAt   time.Time `json:"completed_at"`
	CompletedByID uint      `json:"completed_by"`
}
