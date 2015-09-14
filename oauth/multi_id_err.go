package oauth

import "fmt"

type multiIDErr struct {
	errs []idErr
}

func (e *multiIDErr) addError(err idErr) {
	if e.errs == nil {
		e.errs = []idErr{}
	}
	e.errs = append(e.errs, err)
}

func (e multiIDErr) errors() []idErr {
	return e.errs
}

func (e multiIDErr) Error() string {
	errorMessage := "multiple errors:"
	if e.errs != nil {
		for _, err := range e.errs {
			errorMessage = fmt.Sprintf("%s {%s}", errorMessage, err)
		}
	}
	return errorMessage
}
