package wundergo

import "time"

// Reminder contains information about a task reminder.
type Reminder struct {
	ID        uint      `json:"id"`
	Date      string    `json:"date"`
	TaskID    uint      `json:"task_id"`
	Revision  uint      `json:"revision"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
