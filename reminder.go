package wl

import "time"

// Reminder contains information about a task reminder.
type Reminder struct {
	ID        uint      `json:"id" yaml:"id"`
	Date      string    `json:"date" yaml:"date"`
	TaskID    uint      `json:"task_id" yaml:"task_id"`
	Revision  uint      `json:"revision" yaml:"revision"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}
