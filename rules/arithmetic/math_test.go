package arithmetic

import (
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
		data      map[string]*LeafNode
		mapResult map[string]float64
		wantErr   bool
	}{
		{
			name:    "test c=ac.a+b",
			wantErr: false,
			data: map[string]*LeafNode{
				"default.c": {
					Right: &LeafNode{
						Value: "default.a",
						Type:  ValueType,
					},
					Left: &LeafNode{
						Type:  ValueType,
						Value: "default.b",
					},
					Operator: Add,
					Type:     OperatorType,
				},
			},
			mapResult: map[string]float64{
				"default.c": 4,
			},
		},
		{
			name:    "test d=a-b",
			wantErr: false,
			data: map[string]*LeafNode{
				"default.d": {
					Right: &LeafNode{
						Value: "default.a",
						Type:  ValueType,
					},
					Left: &LeafNode{
						Type:  ValueType,
						Value: "default.b",
					},
					Operator: Subtract,
					Type:     OperatorType,
				},
			},
			mapResult: map[string]float64{
				"default.d": 0,
			},
		},
		{
			name:    "test e=a*b",
			wantErr: false,
			data: map[string]*LeafNode{
				"default.e": {
					Right: &LeafNode{
						Value: "default.a",
						Type:  ValueType,
					},
					Left: &LeafNode{
						Type:  ValueType,
						Value: "default.b",
					},
					Operator: Multiplication,
					Type:     OperatorType,
				},
			},
			mapResult: map[string]float64{
				"default.e": 4,
			},
		},
		{
			name:    "test f=a/b",
			wantErr: false,
			data: map[string]*LeafNode{
				"default.f": {
					Right: &LeafNode{
						Value: "default.a",
						Type:  ValueType,
					},
					Left: &LeafNode{
						Type:  ValueType,
						Value: "default.b",
					},
					Operator: Division,
					Type:     OperatorType,
				},
			},
			mapResult: map[string]float64{
				"default.f": 1,
			},
		},
		{
			name:    "test g=a+10",
			wantErr: false,
			data: map[string]*LeafNode{
				"default.g": {
					Right: &LeafNode{
						Value: "default.a",
						Type:  ValueType,
					},
					Left: &LeafNode{
						Type:  NumberType,
						Value: "10",
					},
					Operator: Add,
					Type:     OperatorType,
				},
			},
			mapResult: map[string]float64{
				"default.g": 12,
			},
		},
		{
			name:    "test y=(a+b)*(a+x)",
			wantErr: false,
			data: map[string]*LeafNode{
				"default.y": {
					Right: &LeafNode{
						Type:     OperatorType,
						Operator: Add,
						Left: &LeafNode{
							Value: "default.a",
							Type:  ValueType,
						},
						Right: &LeafNode{
							Type:  ValueType,
							Value: "default.b",
						},
					},
					Left: &LeafNode{
						Type:     OperatorType,
						Operator: Add,
						Left: &LeafNode{
							Value: "default.a",
							Type:  ValueType,
						},
						Right: &LeafNode{
							Type:  ValueType,
							Value: "default.x",
						},
					},
					Operator: Multiplication,
					Type:     OperatorType,
				},
			},
			mapResult: map[string]float64{
				"default.y": 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := m.Run(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			for key, result := range tt.mapResult {
				q, err := m.dataQuery.Query(key)
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
