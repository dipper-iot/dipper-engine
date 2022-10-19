package daq

import "testing"

var (
	daqData = NewDaq(map[string]interface{}{
		"data1": map[string]interface{}{
			"a": 1,
			"b": 2,
		},
	})
)

func TestDaq_Update_Query(t *testing.T) {

	err := daqData.Update("data1.c", 10)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	q, err := daqData.Query("data1.c")
	if err != nil {
		t.Errorf("Query() error = %v", err)
	}

	data, err := q.Number()
	if err != nil {
		t.Errorf("Number() error = %v", err)
	}
	if data != 10 {
		t.Errorf("Not Match Data")
	}

	q, err = daqData.Query("data1")
	if err != nil {
		t.Errorf("Query() error = %v", err)
	}

	o, err := q.Object()
	if err != nil {
		t.Errorf("Object() error = %v", err)
	}
	err = o.Create("data", map[string]interface{}{
		"c": 1,
	})
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	q, err = daqData.Query("data1.data.c")
	if err != nil {
		t.Errorf("Query() error = %v", err)
	}

	data, err = q.Number()
	if err != nil {
		t.Errorf("Number() error = %v", err)
	}
	if data != 1 {
		t.Errorf("Not Match Data")
	}
}
