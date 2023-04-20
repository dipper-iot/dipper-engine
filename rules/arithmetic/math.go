package arithmetic

import (
	"github.com/Knetic/govaluate"
	"github.com/dipper-iot/dipper-engine/core/daq"
	"github.com/dipper-iot/dipper-engine/pkg/expression"
	"github.com/dipper-iot/dipper-engine/pkg/util"
)

type Math struct {
	mainBranch string
	data       map[string]interface{}
}

func NewMath(mainBranch string, dataQuery map[string]interface{}) *Math {
	return &Math{
		mainBranch: mainBranch,
		data:       util.DataToValue(dataQuery, mainBranch),
	}
}

func (m Math) Run(expressionStr string, keyResult string) error {

	exp, err := govaluate.NewEvaluableExpressionWithFunctions(expressionStr, expression.ExpressionFunctions)
	if err != nil {
		return err
	}

	result, err := exp.Evaluate(m.data)
	if err != nil {
		return err
	}

	queryData := daq.NewDaq(m.data)
	err = queryData.Update(keyResult, result)
	if err != nil {
		return err
	}
	m.data = queryData.Data()

	return nil
}

func (m Math) Data() map[string]interface{} {
	return util.ValueToData(m.data, m.mainBranch)
}
