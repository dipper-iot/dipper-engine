package daq

type ObjectMap struct {
	data  map[string]interface{}
	query *Query
}

func (o ObjectMap) Has(name string) bool {
	_, ok := o.data[name]
	return ok
}

func (o ObjectMap) Create(name string, data interface{}) error {
	o.data[name] = data
	return o.query.updateField(o.query.index, o.data)
}

func (o ObjectMap) Delete(name string) error {
	delete(o.data, name)
	return o.query.updateField(o.query.index, o.data)
}

func (o ObjectMap) CreatePath(location string, data interface{}) error {
	paths := getPath(location)
	return o.createPath(paths, data)
}

func (o ObjectMap) createPath(paths []string, data interface{}) error {

	size := len(paths)
	if size == 1 {
		return o.Create(paths[0], data)
	}

	q, err := o.Query(paths[0]).Object()
	if err != nil {
		return err
	}
	if q.Has(paths[1]) {
		if q.Query(paths[1]).Type() == Object {
			return o.createPath(paths[1:], data)
		}
		return PathExistsNotObject
	}

	return q.createPath(paths[1:], data)
}

func (o ObjectMap) Query(location string) *Query {
	paths := getPath(location)
	return NewQuery(
		o.data,
		append(o.query.paths, paths...),
		o.query.index,
		func(index int, data interface{}) error {
			return o.query.Update(data)
		})
}
