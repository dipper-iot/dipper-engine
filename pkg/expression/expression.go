package expression

import "github.com/Knetic/govaluate"

var ExpressionFunctions = map[string]govaluate.ExpressionFunction{
	"strlen": func(args ...interface{}) (interface{}, error) {
		length := len(args[0].(string))
		return (float64)(length), nil
	},
	"sFromObj": FromObject,
	"nFromObj": NumberFromObject,
}
