package wundergo

import "time"

// TaskComment represents information about a comment on a Task.
type TaskComment struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	Revision  uint      `json:"revision"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
