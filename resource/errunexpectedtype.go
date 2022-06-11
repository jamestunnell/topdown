package resource

import "fmt"

type ErrUnexpectedType struct {
	Actual       any
	ExpectedType string
}

func NewErrUnexpectedType(actual any, expectedType string) *ErrUnexpectedType {
	return &ErrUnexpectedType{Actual: actual, ExpectedType: expectedType}
}

func (err *ErrUnexpectedType) Error() string {
	return fmt.Sprintf("value %v does not have expected type %s", err.Actual, err.ExpectedType)
}
