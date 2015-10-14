package wl

import "time"

// Task contains information about tasks.
// Tasks are children of lists.
type Task struct {
	ID              uint      `json:"id" yaml:"id"`
	AssigneeID      uint      `json:"assignee_id" yaml:"assignee_id"`
	AssignerID      uint      `json:"assigner_id" yaml:"assigner_id"`
	CreatedAt       time.Time `json:"created_at" yaml:"created_at"`
	CreatedByID     uint      `json:"created_by_id" yaml:"created_by_id"`
	DueDate         time.Time `json:"due_date" yaml:"due_date"`
	ListID          uint      `json:"list_id" yaml:"list_id"`
	Revision        uint      `json:"revision" yaml:"revision"`
	Starred         bool      `json:"starred" yaml:"starred"`
	Title           string    `json:"title" yaml:"title"`
	Completed       bool      `json:"completed" yaml:"completed"`
	CompletedAt     time.Time `json:"completed_at" yaml:"completed_at"`
	CompletedByID   uint      `json:"completed_by" yaml:"completed_by"`
	RecurrenceType  string    `json:"recurrence_type" yaml:"recurrence_type"`
	RecurrenceCount uint      `json:"recurrence_count" yaml:"recurrence_count"`
}
