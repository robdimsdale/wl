package wundergo

type Position struct {
	ID       uint   `json:"id"`
	Values   []uint `json:"values"`
	Revision uint   `json:"revision"`
}
