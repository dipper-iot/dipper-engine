package arithmetic

import (
	"github.com/dipper-iot/dipper-engine/core/daq"
	"github.com/dipper-iot/dipper-engine/pkg/util"
	"testing"
)

func TestMath_Run(t *testing.T) {
	m := NewMath("default", map[string]interface{}{
		"default": map[string]interface{}{
			"a": 2,
			"b": 2,
			"x": 3,
			"ac": map[string]interface{}{
				"a": 2,
			},
		},
	})

	tests := []struct {
		name      string
		exp       string
		keyResult string
		mapResult map[string]float64
		wantErr   bool
	}{
		{
			name:      "test c=nFromObj(ac,'a')+b",
			wantErr:   false,
			exp:       "nFromObj(ac,'a')+b",
			keyResult: "c",
			mapResult: map[string]float64{
				"default.c": 4,
			},
		},
		{
			name:      "test d=a-b",
			wantErr:   false,
			exp:       "a-b",
			keyResult: "d",
			mapResult: map[string]float64{
				"default.d": 0,
			},
		},
		{
			name:      "test e=a*b",
			wantErr:   false,
			exp:       "a*b",
			keyResult: "e",
			mapResult: map[string]float64{
				"default.e": 4,
			},
		},
		{
			name:      "test g=a+10",
			wantErr:   false,
			exp:       "a+10",
			keyResult: "g",
			mapResult: map[string]float64{
				"default.g": 12,
			},
		},
		{
			name:      "test y=(a+b)*(a+x)",
			wantErr:   false,
			exp:       "(a+b)*(a+x)",
			keyResult: "y",
			mapResult: map[string]float64{
				"default.y": 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := m.Run(tt.exp, tt.keyResult); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for key, result := range tt.mapResult {
				dataQuery := daq.NewDaq(util.ValueToData(m.data, "default"))
				q, err := dataQuery.Query(key)
				if err != nil {
					t.Errorf("Query() error = %v", err)
					continue
				}
				val, err := q.Number()
				if err != nil {
					t.Errorf("Number() error = %v", err)
					continue
				}
				if val != result {
					t.Errorf("Result Not Match: actual = %f and want = %f", val, result)
				}
			}
		})
	}
}
