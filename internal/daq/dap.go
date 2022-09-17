package daq

type Daq struct {
	data map[string]interface{}
}

func NewDaq(data map[string]interface{}) *Daq {
	return &Daq{
		data: data,
	}
}

func (d Daq) Update(location string, data interface{}) error {
	q, err := d.Query(location)
	if err != nil {
		return err
	}
	return q.Update(data)
}

func (d Daq) Query(location string) (*Query, error) {
	paths := getPath(location)
	q := NewQuery(d.data, paths, 0, func(index int, data interface{}) error {
		path := paths[index]
		d.data[path] = data
		return nil
	})

	return q.getDataIndex(0)
}

func (d Daq) Delete(location string) error {
	paths := getPath(location)
	length := len(paths)
	q := NewQuery(d.data, paths[0:length-2], 0, func(index int, data interface{}) error {
		path := paths[index]
		d.data[path] = data
		return nil
	})

	o, err := q.Object()
	if err != nil {
		return err
	}
	return o.Delete(paths[length-1])
}

func (d *Daq) Clone() *Daq {
	return &Daq{
		data: d.data,
	}
}
