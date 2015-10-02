package wl

import "time"

// User contains information about a User.
type User struct {
	ID        uint      `json:"id" yaml:"id"`
	Name      string    `json:"name" yaml:"name"`
	Email     string    `json:"email" yaml:"email"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
	Revision  uint      `json:"revision" yaml:"revision"`
}
