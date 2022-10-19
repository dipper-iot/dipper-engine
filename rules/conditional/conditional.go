package conditional

import (
	"fmt"
	"github.com/dipper-iot/dipper-engine/core/daq"
	"strconv"
)

type Conditional struct {
	mainBranch string
	dataQuery  *daq.Daq
}

func NewConditional(mainBranch string, dataQuery map[string]interface{}) *Conditional {
	return &Conditional{
		mainBranch: mainBranch,
		dataQuery:  daq.NewDaq(dataQuery),
	}
}

func (r Conditional) Run(leafNode *LeafNode, setResultTo string) (res bool, err error) {

	result, err := r.compare(leafNode)
	if err != nil {
		return false, err
	}

	set := 0
	if result.Result {
		set = 1
	}

	if setResultTo != "" {
		err = r.dataQuery.Update(setResultTo, set)
		if err != nil {
			return false, err
		}
	}

	return result.Result, err
}

func (r Conditional) compare(leafNode *LeafNode) (result *Result, err error) {
	if leafNode.Type == NoneType && leafNode.Left == nil && leafNode.Right == nil {
		return
	}
	result = &Result{
		Result:  false,
		End:     false,
		IsValue: true,
	}

	if leafNode.Type == ValueType {
		var query *daq.Query
		query, err = r.dataQuery.Query(leafNode.Value)
		if err != nil {
			return
		}
		typeData, e := query.QueryTypeItem()
		if e != nil {
			err = e
			return
		}
		switch typeData {
		case daq.Number:
			{
				result.Type = typeData
				result.Value, err = query.Number()
				return
			}
		case daq.String:
			{
				result.Type = typeData
				result.Value, err = query.String()
				return
			}
		case daq.Array:
			{
				result.Type = typeData
				result.Value, err = query.ArrayIndex()
				return
			}
		case daq.Object:
			{
				err = fmt.Errorf("%s value is object", leafNode.Value)
				return
			}
		}

		return
	}
	if leafNode.Type == NumberType {
		result.Type = daq.Number
		result.Value, err = strconv.ParseFloat(leafNode.Value, 10)
		return
	}
	if leafNode.Type != OperatorType {
		err = fmt.Errorf("not found Type: %s", leafNode.Type)
		return
	}
	var left, right *Result
	if leafNode.Left != nil {
		left, err = r.compare(leafNode.Left)
		if err != nil {
			return
		}
	}
	if leafNode.Right != nil {
		right, err = r.compare(leafNode.Right)
		if err != nil {
			return
		}
	}
	var (
		dataCompare *ResultCompare
	)
	dataCompare, err = r.getCompareResult(left, right)
	if err != nil {
		return
	}

	result.End = true

	if dataCompare.IsResult {
		result.Result, err = r.getResultLogic(dataCompare, leafNode.Operator)
		return
	}

	result.Result, err = r.getResult(dataCompare, leafNode.Operator)

	return
}

func (r Conditional) Data() map[string]interface{} {
	return r.dataQuery.Data()
}

func (r Conditional) getCompareResult(left, right *Result) (result *ResultCompare, err error) {
	result = &ResultCompare{
		IsString: left.Type == daq.String,
	}
	if left.End || right.End {
		if (left.End && !right.End) || (!left.End && right.End) {
			err = fmt.Errorf("not match compare two type other")
			return
		}
		result.IsResult = true
		result.LeftResult = left.Result
		result.RightResult = right.Result
		return
	}

	if left.Type != right.Type {
		err = fmt.Errorf("not match type compare")
		return
	}
	if left.Type == daq.String {
		result.LeftString = left.Value.(string)
		result.RightString = right.Value.(string)
		return
	}
	if left.Type == daq.Number {
		result.LeftNumber = left.Value.(float64)
		result.RightNumber = right.Value.(float64)
		return
	}
	err = fmt.Errorf("type have not support compare")
	return
}

func (r Conditional) getResult(data *ResultCompare, operator Operator) (result bool, err error) {
	if data.IsString {
		switch operator {
		case Equal:
			{
				result = data.LeftString == data.RightString
				break
			}
		case Difference:
			{
				result = data.LeftString != data.RightString
				break
			}
		case LessThan, LessThanOrEqual, GreaterThanOrEqual, GreaterThan:
			{
				err = fmt.Errorf("type have not support operator string")
				break
			}
		}
		return
	}

	switch operator {
	case Equal:
		{
			result = data.LeftNumber == data.RightNumber
			return
		}
	case Difference:
		{
			result = data.LeftNumber != data.RightNumber
			return
		}
	case LessThan:
		{
			result = data.LeftNumber < data.RightNumber
			return
		}
	case LessThanOrEqual:
		{
			result = data.LeftNumber >= data.RightNumber
			return
		}
	case GreaterThanOrEqual:
		{
			result = data.LeftNumber <= data.RightNumber
			return
		}
	case GreaterThan:
		{
			result = data.LeftNumber > data.RightNumber
			return
		}
	}

	err = fmt.Errorf("type have not support compare")
	return
}

func (r Conditional) getResultLogic(data *ResultCompare, operator Operator) (result bool, err error) {
	if data.IsResult {
		switch operator {
		case And:
			{
				result = data.LeftResult && data.RightResult
				return
			}
		case Or:
			{
				result = data.LeftResult || data.RightResult
				return
			}
		}
	}

	err = fmt.Errorf("type have not support compare logic")
	return
}
