package daq

import "errors"

var (
	NotFoundPath         = errors.New("NotFoundPath")
	NotConvertTypeNumber = errors.New("NotConvertTypeNumber")
	NotConvertTypeString = errors.New("NotConvertTypeString")
	NotConvertTypeArray  = errors.New("NotConvertTypeArray")
	NotConvertTypeObject = errors.New("NotConvertTypeObject")
	NotArrayIndex        = errors.New("NotArrayIndex")
	PathExistsNotObject  = errors.New("PathExistsNotObject")
)
