package arithmetic

type DataType string
type Operator string

const (
	ValueType    DataType = "val"
	NumberType   DataType = "number"
	OperatorType DataType = "operator"
	NoneType     DataType = "none"
)

const (
	Add            Operator = "add"
	Subtract       Operator = "subtract"
	Multiplication Operator = "multiplication"
	Division       Operator = "division"
)

type LeafNode struct {
	Left     *LeafNode `json:"left"`
	Right    *LeafNode `json:"right"`
	Type     DataType  `json:"type"`
	Operator Operator  `json:"operator"`
	Value    string    `json:"value"`
}

type Option struct {
	List        map[string]*LeafNode `json:"list"`
	NextError   string               `json:"next_error"`
	NextSuccess string               `json:"next_success"`
	Debug       bool                 `json:"debug"`
}
