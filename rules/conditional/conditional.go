package conditional

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/dipper-iot/dipper-engine/core/daq"
	"github.com/dipper-iot/dipper-engine/pkg/expression"
	"github.com/dipper-iot/dipper-engine/pkg/util"
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

func (r Conditional) Run(conditional string, setResultTo string) (res bool, err error) {

	exp, err := govaluate.NewEvaluableExpressionWithFunctions(conditional, expression.ExpressionFunctions)
	if err != nil {
		return false, err
	}

	dataExp := util.DataToValue(r.dataQuery.Data(), r.mainBranch)
	raw, err := exp.Evaluate(dataExp)
	if err != nil {
		return false, err
	}
	result, ok := raw.(bool)
	if !ok {
		return false, errors.New("result not bool")
	}
	set := 0
	if result {
		set = 1
	}

	if setResultTo != "" {
		mainBranch := r.mainBranch
		if mainBranch == "" {
			mainBranch = "default"
		}
		err = r.dataQuery.Update(fmt.Sprintf("%s.%s", mainBranch, setResultTo), set)
		if err != nil {
			return false, err
		}
	}

	return result, err
}

func (r Conditional) Data() map[string]interface{} {
	return util.ValueToData(r.dataQuery.Data(), r.mainBranch)
}
