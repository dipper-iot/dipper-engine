package expression

import (
	"errors"
	"github.com/dipper-iot/dipper-engine/core/daq"
)

var FromObject = func(args ...interface{}) (interface{}, error) {
	dataMap, ok := args[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("not a object")
	}
	key, ok := args[1].(string)
	if !ok {
		return nil, errors.New("have key get")
	}
	queryData := daq.NewDaq(dataMap)
	res, err := queryData.Query(key)
	if err != nil {
		return nil, err
	}
	return res.Interface()
}

var NumberFromObject = func(args ...interface{}) (interface{}, error) {
	dataMap, ok := args[0].(map[string]interface{})
	if !ok {
		return nil, errors.New("not a object")
	}
	key, ok := args[1].(string)
	if !ok {
		return nil, errors.New("have key get")
	}
	queryData := daq.NewDaq(dataMap)
	res, err := queryData.Query(key)
	if err != nil {
		return nil, err
	}
	return res.Number()
}
