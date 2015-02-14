package wundergo

import "time"

// User contains information about a User.
type User struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Revision   uint      `json:"revision"`
	TypeString string    `json:"type"`
}
