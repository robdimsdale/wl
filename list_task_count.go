package wundergo

// ListTaskCount contains information about the number and type of tasks
// a List contains.
type ListTaskCount struct {
	ID               uint   `json:"id"`
	CompletedCount   uint   `json:"completed_count"`
	UncompletedCount uint   `json:"uncompleted_count"`
	TypeString       string `json:"type"`
}
