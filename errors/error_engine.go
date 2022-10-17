package errors

import "fmt"

type ErrorEngine struct {
	Id          string
	SessionId   uint64
	FromEngine  string
	Message     string
	ErrorDetail error
	Code        int
}

func (e ErrorEngine) Error() string {
	return fmt.Sprintf("%v", e)
}
