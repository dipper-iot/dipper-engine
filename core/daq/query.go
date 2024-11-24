package daq

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type CallUpdate = func(index int, data interface{}) error

type Query struct {
	data           interface{}
	paths          []string
	index          int
	update         CallUpdate
	dataObject     map[string]interface{}
	dataNumber     float64
	dataString     string
	dataArray      []interface{}
	dataIndexArray int
}

func NewQuery(data interface{}, paths []string, index int, update CallUpdate) *Query {
	return &Query{data: data, paths: paths, index: index, update: update}
}

func (q *Query) getDataIndex(index int) (*Query, error) {
	length := len(q.paths) - 1
	if index == length {
		return q, nil
	}

	typeData := q.Type()

	switch typeData {
	case Number:
	case String:
	case Array:
		{
			return q, fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
		}
	default:

	}

	path := q.paths[index]
	typePath := getPathType(path)

	// array
	if typePath == ArrayData {
		dataObject, err := q.ArrayIndex()
		if err != nil {
			return q, err
		}

		nextQuery := NewQuery(
			dataObject,
			q.paths,
			q.index+1,
			q.updateField,
		)
		return nextQuery.getDataIndex(index + 1)
	}

	// object
	dataObject, err := q.object()
	if err != nil {
		return q, err
	}
	newData, ok := dataObject[path]
	if !ok {
		return q, fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
	}

	nextQuery := NewQuery(
		newData,
		q.paths,
		q.index+1,
		q.updateField,
	)
	return nextQuery.getDataIndex(index + 1)
}

func (q *Query) updateField(index int, data interface{}) error {

	var output interface{}
	switch q.Type() {
	case Object:
		_, err := q.object()
		if err != nil {
			return err
		}
		path := q.paths[index]
		q.dataObject[path] = data
		output = q.dataObject
		break
	case Array:
		_, err := q.Array()
		if err != nil {
			return err
		}
		q.dataArray[q.dataIndexArray] = data
		output = q.dataArray
		break
	}

	if index == 0 {
		return nil
	}

	return q.update(index, output)
}

func (q *Query) Update(data interface{}) error {
	var output interface{}
	switch q.Type() {
	case Object:
		_, err := q.object()
		if err != nil {
			return err
		}
		path := q.paths[q.index]
		q.dataObject[path] = data
		output = q.dataObject
		break
	case Array:
		_, err := q.Array()
		if err != nil {
			return err
		}
		q.dataArray[q.dataIndexArray] = data
		output = q.dataArray
		break
	}
	if q.index == 0 {
		return nil
	}

	return q.update(q.index-1, output)
}

func (q *Query) Type() TypeData {

	v := reflect.ValueOf(q.data)
	if v.CanFloat() {
		return Number
	}

	vArray := reflect.ValueOf([]interface{}{})
	if v.CanConvert(vArray.Type()) {
		return Array
	}

	vObject := reflect.ValueOf(map[string]interface{}{})
	if v.CanConvert(vObject.Type()) {
		return Object
	}

	return String
}

func (q *Query) QueryTypeItem() (TypeData, error) {
	if q.index > len(q.paths)-1 {
		return 0, fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
	}
	if q.Type() == Object {
		name := q.paths[q.index]
		data, err := q.object()
		if err != nil {
			return 0, err
		}
		n, ok := data[name]
		if !ok {
			return 0, fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
		}
		return toType(n), nil
	}

	return toType(q.data), nil
}

func toType(data interface{}) TypeData {

	v := reflect.ValueOf(data)
	if v.CanInt() {
		return Number
	}
	tMap := reflect.ValueOf(map[string]interface{}{})
	if v.CanConvert(tMap.Type()) {
		return Object
	}
	return String
}

func (q *Query) Number() (float64, error) {
	if q.dataNumber != 0 {
		return q.dataNumber, nil
	}
	if q.index > len(q.paths)-1 {
		return 0, fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
	}
	if q.Type() == Object {

		name := q.paths[q.index]
		data, err := q.object()
		if err != nil {
			return 0, err
		}
		n, ok := data[name]
		if !ok {
			return 0, fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
		}
		switch v := n.(type) {
		case string:
			s, err := strconv.ParseFloat(v, 64)
			return s, err
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case int:
			return float64(v), nil
		default:
		}

		v := reflect.ValueOf(n)
		if v.CanInt() {
			v2 := v.Int()
			return float64(v2), nil
		}
		if !v.CanFloat() {
			return 0, NotConvertTypeNumber
		}

		q.dataNumber = v.Float()
		return q.dataNumber, nil

		return 0, nil
	}
	v := reflect.ValueOf(q.data)
	if !v.CanFloat() {
		return 0, NotConvertTypeNumber
	}

	q.dataNumber = v.Float()
	return q.dataNumber, nil
}

func (q *Query) String() (string, error) {
	if q.dataString != "" {
		return q.dataString, nil
	}
	if q.index > len(q.paths)-1 {
		return "", fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
	}
	if q.Type() == Object {
		name := q.paths[q.index]
		data, err := q.object()
		if err != nil {
			return "", err
		}
		n, ok := data[name]
		if !ok {
			return "", fmt.Errorf("%s: %s", NotFoundPath, strings.Join(q.paths, "."))
		}
		switch v := n.(type) {
		case string:
			return v, nil
		case float64, float32:
			return fmt.Sprintf("%f", v), nil
		case int, int64:
			return fmt.Sprintf("%d", v), nil
		default:
		}
		v := reflect.ValueOf(n)
		q.dataString = v.String()
		return q.dataString, nil
	}
	v := reflect.ValueOf(q.data)
	q.dataString = v.String()
	return q.dataString, nil
}

func (q *Query) Interface() (interface{}, error) {

	mapData, err := q.Object()
	if err != nil {
		return nil, err
	}
	name := q.paths[q.index]
	if !mapData.Has(name) {
		return nil, errors.New("not found " + name)
	}
	return mapData.data[name], nil
}

func (q *Query) object() (map[string]interface{}, error) {
	if q.dataObject != nil {
		return q.dataObject, nil
	}
	v := reflect.ValueOf(q.data)

	q.dataObject = make(map[string]interface{})
	out := reflect.ValueOf(q.dataObject)

	if !v.CanConvert(out.Type()) {
		return nil, NotConvertTypeObject
	}
	val := v.Convert(out.Type())

	if !val.IsValid() {
		return nil, NotConvertTypeObject
	}
	if val.IsNil() || out.IsNil() {
		return nil, NotConvertTypeObject
	}
	i := val.Interface()
	q.dataObject = i.(map[string]interface{})

	return q.dataObject, nil
}

func (q *Query) Object() (*ObjectMap, error) {
	// object
	dataObject, err := q.object()
	if err != nil {
		return nil, err
	}

	return &ObjectMap{
		data:  dataObject,
		query: q,
	}, nil
}

func (q *Query) Array() ([]interface{}, error) {
	if q.dataArray != nil {
		return q.dataArray, nil
	}
	v := reflect.ValueOf(q.data)
	length := v.Len()
	for i := 0; i < length; i++ {
		val := v.Index(i)
		if val.CanInterface() {
			return nil, NotConvertTypeArray
		}
		q.dataArray = append(q.dataArray, val.Interface())
	}
	return q.dataArray, nil
}

func (q *Query) ArrayIndex() (interface{}, error) {
	data, err := q.Array()
	if err != nil {
		return nil, err
	}
	index, err := getIndexArray(q.paths[q.index])
	if err != nil {
		return nil, err
	}
	q.dataIndexArray = index
	if q.dataIndexArray > len(q.dataArray) {
		return nil, NotArrayIndex
	}
	return data[q.dataIndexArray], nil
}
