package wundergo

import "time"

// Folder contains information about a folder.
// A list may or may not belong to a folder.
type Folder struct {
	ID                 uint      `json:"id"`
	Title              string    `json:"title"`
	ListIDs            []uint    `json:"list_ids"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedByRequestID uint      `json:"created_by_request_id"`
	UpdatedAt          time.Time `json:"updated_at"`
	TypeString         string    `json:"type"`
	Revision           uint      `json:"revision"`
}
