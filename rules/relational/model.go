package relational

import "github.com/dipper-iot/dipper-engine/internal/daq"

type DataType string
type Operator string

const (
	ValueType    DataType = "val"
	NumberType   DataType = "number"
	OperatorType DataType = "operator"
	NoneType     DataType = "none"
)

const (
	Equal              Operator = "=="
	Difference         Operator = "<>"
	LessThanOrEqual    Operator = "<="
	LessThan           Operator = "<"
	GreaterThan        Operator = ">"
	GreaterThanOrEqual Operator = ">="
	And                Operator = "&&"
	Or                 Operator = "||"
)

type LeafNode struct {
	Left     *LeafNode `json:"left"`
	Right    *LeafNode `json:"right"`
	Type     DataType  `json:"type"`
	Operator Operator  `json:"operator"`
	Value    string    `json:"data"`
}

type Option struct {
	Operator         *LeafNode `json:"operator"`
	SetParamResultTo string    `json:"set_param_result_to"`
	NextError        string    `json:"next_error"`
	NextTrue         string    `json:"next_true"`
	NextFalse        string    `json:"next_false"`
	Debug            bool      `json:"debug"`
}

type Result struct {
	End     bool
	Result  bool
	IsValue bool
	Type    daq.TypeData
	Value   interface{}
}

type ResultCompare struct {
	IsString bool

	LeftNumber  float64
	RightNumber float64
	LeftString  string
	RightString string

	IsResult    bool
	LeftResult  bool
	RightResult bool
}
