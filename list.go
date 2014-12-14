package wundergo

type List struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	CreatedAt  string `json:"created_at"`
	ListType   string `json:"list_type"`
	Revision   int    `json:"revision"`
	TypeString string `json:"type"`
	Public     bool   `json:"public"`
}

type ListTaskCount struct {
	ID               uint   `json:"id"`
	CompletedCount   int    `json:"completed_count"`
	UncompletedCount int    `json:"uncompleted_count"`
	TypeString       string `json:"type"`
}
