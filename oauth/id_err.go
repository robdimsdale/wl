package oauth

import "fmt"

type idErr struct {
	idType string
	id     uint
	err    error
}

func (e idErr) Error() string {
	return fmt.Sprintf("%s id: %d, err: %v", e.idType, e.id, e.err)
}
