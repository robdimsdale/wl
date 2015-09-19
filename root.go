package wundergo

// Root contains information about the root of the object hierarchy.
type Root struct {
	ID       uint `json:"id"`
	Revision uint `json:"revision"`
	UserID   uint `json:"user_id"`
}
