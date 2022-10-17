package arithmetic

import (
	"fmt"
	"github.com/dipper-iot/dipper-engine/internal/daq"
	"strconv"
)

type Math struct {
	mainBranch string
	dataQuery  *daq.Daq
}

func NewMath(mainBranch string, dataQuery map[string]interface{}) *Math {
	return &Math{
		mainBranch: mainBranch,
		dataQuery:  daq.NewDaq(dataQuery),
	}
}

func (m Math) Run(leafNodes map[string]*LeafNode) error {

	for keyResult, node := range leafNodes {
		result, err := m.calculator(node)
		if err != nil {
			return err
		}
		err = m.dataQuery.Update(keyResult, result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Math) calculator(leafNode *LeafNode) (result float64, err error) {
	if leafNode.Type == NoneType && leafNode.Left == nil && leafNode.Right == nil {
		return
	}
	if leafNode.Type == ValueType {
		var query *daq.Query
		query, err = m.dataQuery.Query(leafNode.Value)
		if err != nil {
			return
		}

		result, err = query.Number()
		return
	}
	if leafNode.Type == NumberType {
		result, err = strconv.ParseFloat(leafNode.Value, 10)
		return
	}
	if leafNode.Type != OperatorType {
		err = fmt.Errorf("not found Type: %s", leafNode.Type)
		return
	}
	var left, right float64
	if leafNode.Left != nil {
		left, err = m.calculator(leafNode.Left)
		if err != nil {
			return
		}
	}
	if leafNode.Right != nil {
		right, err = m.calculator(leafNode.Right)
		if err != nil {
			return
		}
	}

	switch leafNode.Operator {
	case Add:
		{
			result = left + right
			break
		}
	case Subtract:
		{
			result = left - right
			break
		}
	case Multiplication:
		{
			result = left * right
			break
		}
	case Division:
		{
			if right == 0 {
				err = fmt.Errorf("division by zero")
				return
			}
			result = left / right
			break
		}
	}
	return
}

func (m Math) Data() map[string]interface{} {
	return m.dataQuery.Data()
}
