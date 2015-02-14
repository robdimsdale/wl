package wundergo

import "time"

// Task contains information about tasks.
// Tasks are children of lists.
type Task struct {
	ID              uint      `json:"id"`
	AssigneeID      uint      `json:"assignee_id"`
	AssignerID      uint      `json:"assigner_id"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedByID     uint      `json:"created_by_id"`
	DueDate         string    `json:"due_date"`
	ListID          uint      `json:"list_id"`
	Revision        uint      `json:"revision"`
	Starred         bool      `json:"starred"`
	Title           string    `json:"title"`
	Completed       bool      `json:"completed"`
	CompletedAt     time.Time `json:"completed_at"`
	CompletedByID   uint      `json:"completed_by"`
	RecurrenceType  string    `json:"recurrence_type"`
	RecurrenceCount uint      `json:"recurrence_count"`
}
