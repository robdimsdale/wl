package wl

import "time"

// List contains information about a List.
type List struct {
	ID         uint      `json:"id" yaml:"id"`
	Title      string    `json:"title" yaml:"title"`
	CreatedAt  time.Time `json:"created_at" yaml:"created_at"`
	ListType   string    `json:"list_type" yaml:"list_type"`
	Revision   uint      `json:"revision" yaml:"revision"`
	TypeString string    `json:"type" yaml:"type"`
	Public     bool      `json:"public" yaml:"public"`
}
