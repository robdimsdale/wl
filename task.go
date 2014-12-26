package wundergo

type Task struct {
	ID          uint   `json:"id"`
	AssigneeID  uint   `json:"assignee_id"`
	AssignerID  uint   `json:"assigner_id"`
	CreatedAt   string `json:"created_at"`
	CreatedByID uint   `json:"created_by_id"`
	DueDate     string `json:"dueDate"`
	ListID      uint   `json:"list_id"`
	Revision    int    `json:"revision"`
	Starred     bool   `json:"starred"`
	Title       string `json:"title"`
	Completed   bool   `json:"completed"`
}
