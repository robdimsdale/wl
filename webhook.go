package wl

import "time"

// Webhook contains information about webhooks.
// A webhook sends notifications when a list is updated.
type Webhook struct {
	ID             uint      `json:"id" yaml:"id"`
	ListID         uint      `json:"list_id" yaml:"list_id"`
	MembershipID   uint      `json:"membership_id" yaml:"membership_id"`
	MembershipType string    `json:"membership_type" yaml:"membership_type"`
	URL            string    `json:"url" yaml:"url"`
	ProcessorType  string    `json:"processor_type" yaml:"processor_type"`
	Configuration  string    `json:"configuration" yaml:"configuration"`
	CreatedAt      time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" yaml:"updated_at"`
}
