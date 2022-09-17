package errors

import "fmt"

type ErrorEngine struct {
	Id         string
	RequestId  string
	FromEngine string
	Message    string
	Code       int
}

func (e ErrorEngine) Error() string {
	return fmt.Sprintf("%v", e)
}
