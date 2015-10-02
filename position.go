package wl

// Position contains an ordered list of IDs of Lists, Tasks or Subtasks.
type Position struct {
	ID       uint   `json:"id" yaml:"id"`
	Values   []uint `json:"values" yaml:"values"`
	Revision uint   `json:"revision" yaml:"revision"`
}
