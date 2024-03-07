package core

import "errors"

type DataSource[T any] interface {
	Load(fileName string) (T, error)
}

type Identifiable interface {
	GetId() int
}

type DictionaryData[T Identifiable] struct {
	fileName string
	name     string
	data     map[int]T
}

func NewDictionaryData[T Identifiable](fileName string, name string, list []T) DictionaryData[T] {
	data := make(map[int]T)
	for _, v := range list {
		data[v.GetId()] = v
	}
	return DictionaryData[T]{fileName, name, data}
}

func (d *DictionaryData[T]) Get(idx int) (*T, error) {
	v, ok := d.data[idx]
	if !ok {
		return nil, errors.New("invalid " + d.name + " id")
	}
	return &v, nil
}
