package wundergo

import "time"

// Webhook contains information about webhooks.
// A webhook sends notifications when a list is updated.
type Webhook struct {
	ID             uint      `json:"id"`
	ListID         uint      `json:"list_id"`
	MembershipID   uint      `json:"membership_id"`
	MembershipType string    `json:"membership_type"`
	URL            string    `json:"url"`
	ProcessorType  string    `json:"processor_type"`
	Configuration  string    `json:"configuration"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
