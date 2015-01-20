package wundergo

type Membership struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	ListID uint   `json:"list_id"`
	State  string `json:"state"`
	Owner  bool   `json:"owner"`
	Muted  bool   `json:"muted"`
}
