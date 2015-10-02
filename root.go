package wl

// Root contains information about the root of the object hierarchy.
type Root struct {
	ID       uint `json:"id" yaml:"id"`
	Revision uint `json:"revision" yaml:"revision"`
	UserID   uint `json:"user_id" yaml:"user_id"`
}
