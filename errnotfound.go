package topdown

import "fmt"

type ErrNotFound struct {
	Name string
}

func NewErrNotFound(name string) *ErrNotFound {
	return &ErrNotFound{Name: name}
}

func (err *ErrNotFound) Error() string {
	return fmt.Sprintf("'%s' not found", err.Name)
}
