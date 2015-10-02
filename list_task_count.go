package wl

// ListTaskCount contains information about the number and type of tasks
// a List contains.
type ListTaskCount struct {
	ID               uint   `json:"id" yaml:"id"`
	CompletedCount   uint   `json:"completed_count" yaml:"completed_count"`
	UncompletedCount uint   `json:"uncompleted_count" yaml:"uncompleted_count"`
	TypeString       string `json:"type" yaml:"type"`
}
