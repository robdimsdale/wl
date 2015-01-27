package wundergo

// List contains information about a List.
type List struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	CreatedAt  string `json:"created_at"`
	ListType   string `json:"list_type"`
	Revision   uint   `json:"revision"`
	TypeString string `json:"type"`
	Public     bool   `json:"public"`
}
