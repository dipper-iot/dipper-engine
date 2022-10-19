package errors

import (
	"errors"
)

var (
	ErrorNotFoundQueue    = errors.New("NoT Found Queue")
	ErrorNotFoundProperty = errors.New("NoT Found Property")
)
