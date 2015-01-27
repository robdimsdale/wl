package wundergo

// Position contains an ordered list of IDs of Lists, Tasks or Subtasks.
type Position struct {
	ID       uint   `json:"id"`
	Values   []uint `json:"values"`
	Revision uint   `json:"revision"`
}
