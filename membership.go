package wl

// Membership joins Users and Lists.
type Membership struct {
	ID       uint   `json:"id" yaml:"id"`
	UserID   uint   `json:"user_id" yaml:"user_id"`
	ListID   uint   `json:"list_id" yaml:"list_id"`
	State    string `json:"state" yaml:"state"`
	Owner    bool   `json:"owner" yaml:"owner"`
	Muted    bool   `json:"muted" yaml:"muted"`
	Revision uint   `json:"revision" yaml:"revision"`
}
