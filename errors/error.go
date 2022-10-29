package errors

import (
	"errors"
)

var (
	ErrorNotFoundQueue    = errors.New("NoT Found Queue")
	ErrorNotFoundProperty = errors.New("NoT Found Property")
	ErrorNotControlEngine = errors.New("Not Have Session Control With Rule Infinity")
)
