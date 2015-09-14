package oauth

import "fmt"

type idErr struct {
	id  uint
	err error
}

func (e idErr) Error() string {
	return fmt.Sprintf("id: %d, err: %v", e.id, e.err)
}
