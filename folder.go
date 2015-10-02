package wl

import "time"

// Folder contains information about a folder.
// A list may or may not belong to a folder.
type Folder struct {
	ID                 uint      `json:"id" yaml:"id"`
	Title              string    `json:"title" yaml:"title"`
	ListIDs            []uint    `json:"list_ids" yaml:"list_ids"`
	CreatedAt          time.Time `json:"created_at" yaml:"created_at"`
	CreatedByRequestID string    `json:"created_by_request_id" yaml:"created_by_request_id"`
	UpdatedAt          time.Time `json:"updated_at" yaml:"updated_at"`
	TypeString         string    `json:"type" yaml:"type"`
	Revision           uint      `json:"revision" yaml:"revision"`
}
