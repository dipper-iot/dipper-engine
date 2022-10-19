package _switch

import (
	"fmt"
	"github.com/dipper-iot/dipper-engine/core/daq"
	"reflect"
)

type Switch struct {
	mainBranch string
	dataQuery  *daq.Daq
}

func NewSwitch(mainBranch string, dataQuery map[string]interface{}) *Switch {
	return &Switch{
		mainBranch: mainBranch,
		dataQuery:  daq.NewDaq(dataQuery),
	}
}

func (r Switch) Run(dataMap map[string]string, key string) (result string, err error) {

	var query *daq.Query
	query, err = r.dataQuery.Query(key)
	if err != nil {
		return
	}
	var (
		typData    daq.TypeData
		strData    string
		numberData float64
	)
	typData, err = query.QueryTypeItem()
	if err != nil {
		return
	}
	switch typData {
	case daq.String:
		{
			strData, err = query.String()
			if err != nil {
				return
			}
			break
		}
	case daq.Number:
		{
			strData, err = query.String()
			if err != nil {
				return
			}
			break
		}
	default:
		{
			err = fmt.Errorf("commpare type ")
			return
		}

	}

	for redirect, dataMatch := range dataMap {
		if typData == daq.String && dataMatch == strData {
			result = redirect
			return
		}
		if typData == daq.Number {
			v := reflect.ValueOf(dataMatch)
			if !v.CanFloat() {
				err = fmt.Errorf("not convert: %s", dataMatch)
				return
			}
			if numberData == v.Float() {
				result = redirect
				return
			}
		}
	}

	result = dataMap["default"]
	return
}
