package conditional

import (
	"testing"
)

func TestConditional_Run(t *testing.T) {
	m := NewConditional("default", map[string]interface{}{
		"default": map[string]interface{}{
			"a": 2,
			"b": 2,
			"x": 3,
		},
	})

	tests := []struct {
		name      string
		data      *LeafNode
		setTo     string
		mapResult map[string]bool
		wantErr   bool
	}{
		{
			name:    "test c=a==b",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Value: "default.a",
					Type:  ValueType,
				},
				Left: &LeafNode{
					Type:  ValueType,
					Value: "default.b",
				},
				Operator: Equal,
				Type:     OperatorType,
			},
			setTo: "default.c",
			mapResult: map[string]bool{
				"default.c": true,
			},
		},
		{
			name:    "test d=a<>b",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Value: "default.a",
					Type:  ValueType,
				},
				Left: &LeafNode{
					Type:  ValueType,
					Value: "default.b",
				},
				Operator: Difference,
				Type:     OperatorType,
			},
			setTo: "default.d",
			mapResult: map[string]bool{
				"default.d": false,
			},
		},
		{
			name:    "test e=a>=b",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Value: "default.a",
					Type:  ValueType,
				},
				Left: &LeafNode{
					Type:  ValueType,
					Value: "default.b",
				},
				Operator: GreaterThanOrEqual,
				Type:     OperatorType,
			},
			setTo: "default.e",
			mapResult: map[string]bool{
				"default.e": true,
			},
		},
		{
			name:    "test f=a<=b",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Value: "default.a",
					Type:  ValueType,
				},
				Left: &LeafNode{
					Type:  ValueType,
					Value: "default.b",
				},
				Operator: LessThanOrEqual,
				Type:     OperatorType,
			},
			setTo: "default.f",
			mapResult: map[string]bool{
				"default.f": true,
			},
		},
		{
			name:    "test g=a>10",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Value: "default.a",
					Type:  ValueType,
				},
				Left: &LeafNode{
					Type:  ValueType,
					Value: "default.b",
				},
				Operator: LessThan,
				Type:     OperatorType,
			},
			setTo: "default.g",
			mapResult: map[string]bool{
				"default.g": false,
			},
		},
		{
			name:    "test h=a>10",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Value: "default.a",
					Type:  ValueType,
				},
				Left: &LeafNode{
					Type:  ValueType,
					Value: "default.b",
				},
				Operator: GreaterThan,
				Type:     OperatorType,
			},
			setTo: "default.h",
			mapResult: map[string]bool{
				"default.h": false,
			},
		},
		{
			name:    "test y=(a==b)&&(a>x)",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Type:     OperatorType,
					Operator: Equal,
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
					Operator: GreaterThan,
					Left: &LeafNode{
						Value: "default.a",
						Type:  ValueType,
					},
					Right: &LeafNode{
						Type:  ValueType,
						Value: "default.x",
					},
				},
				Operator: And,
				Type:     OperatorType,
			},
			setTo: "default.y",
			mapResult: map[string]bool{
				"default.y": false,
			},
		},
		{
			name:    "test z=(a==b)||(a>x)",
			wantErr: false,
			data: &LeafNode{
				Right: &LeafNode{
					Type:     OperatorType,
					Operator: Equal,
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
					Operator: GreaterThan,
					Left: &LeafNode{
						Value: "default.a",
						Type:  ValueType,
					},
					Right: &LeafNode{
						Type:  ValueType,
						Value: "default.x",
					},
				},
				Operator: Or,
				Type:     OperatorType,
			},
			setTo: "default.z",
			mapResult: map[string]bool{
				"default.z": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if _, err := m.Run(tt.data, tt.setTo); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
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
				var want float64 = 0
				if result {
					want = 1
				}
				if val != want {
					t.Errorf("Result Not Match: actual = %f and want = %v", val, result)
				}
			}
		})
	}
}
