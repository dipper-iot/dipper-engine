package conditional

import (
	"github.com/dipper-iot/dipper-engine/core/daq"
)

type Option struct {
	Conditional      string `json:"conditional"`
	SetParamResultTo string `json:"set_param_result_to"`
	NextError        string `json:"next_error"`
	NextTrue         string `json:"next_true"`
	NextFalse        string `json:"next_false"`
	Debug            bool   `json:"debug"`
}

type Result struct {
	End     bool
	Result  bool
	IsValue bool
	Type    daq.TypeData
	Value   interface{}
}
