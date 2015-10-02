package wl

import "time"

// Note represents the information about a note.
// Notes are large text blobs, and are children of tasks.
type Note struct {
	ID        uint      `json:"id" yaml:"id"`
	TaskID    uint      `json:"task_id" yaml:"task_id"`
	Content   string    `json:"content" yaml:"content"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
	Revision  uint      `json:"revision" yaml:"revision"`
}
