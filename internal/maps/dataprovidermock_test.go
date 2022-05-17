package maps

import "reflect"

type mockedResourceDataProvider struct {
	actual map[string]interface{}
	old    map[string]interface{}
}

func (r mockedResourceDataProvider) Get(key string) interface{} {
	return r.actual[key]
}

func (r mockedResourceDataProvider) GetOk(key string) (interface{}, bool) {
	v, ok := r.actual[key]
	return v, ok
}

func (r mockedResourceDataProvider) HasChange(key string) bool {
	return !reflect.DeepEqual(r.old[key], r.actual[key])
}

func (r mockedResourceDataProvider) GetChange(key string) (interface{}, interface{}) {
	return r.old[key], r.actual[key]
}
