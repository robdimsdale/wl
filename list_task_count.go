package wundergo

type ListTaskCount struct {
	ID               uint   `json:"id"`
	CompletedCount   int    `json:"completed_count"`
	UncompletedCount int    `json:"uncompleted_count"`
	TypeString       string `json:"type"`
}
