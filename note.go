package wundergo

import "time"

// Note represents the information about a note.
// Notes are large text blobs, and are children of tasks.
type Note struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Revision  uint      `json:"revision"`
}
