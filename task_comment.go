package wl

import "time"

// TaskComment represents information about a comment on a Task.
type TaskComment struct {
	ID        uint      `json:"id" yaml:"id"`
	TaskID    uint      `json:"task_id" yaml:"task_id"`
	Revision  uint      `json:"revision" yaml:"revision"`
	Text      string    `json:"text" yaml:"text"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}
