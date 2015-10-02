package wl

import "time"

// Subtask contains information about a subtask.
// Subtasks are children of tasks.
type Subtask struct {
	ID            uint      `json:"id" yaml:"id"`
	TaskID        uint      `json:"task_id" yaml:"task_id"`
	CreatedAt     time.Time `json:"created_at" yaml:"created_at"`
	CreatedByID   uint      `json:"created_by_id" yaml:"created_by_id"`
	Revision      uint      `json:"revision" yaml:"revision"`
	Title         string    `json:"title" yaml:"title"`
	Completed     bool      `json:"completed" yaml:"completed"`
	CompletedAt   time.Time `json:"completed_at" yaml:"completed_at"`
	CompletedByID uint      `json:"completed_by" yaml:"completed_by"`
}
