package core

import "errors"

type DataSource[T any] interface {
	Load(fileName string) (T, error)
}

type DataSaver[T any] interface {
	Save(data *T, fileName string) error
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

func (d *DictionaryData[T]) SaveTo(saver DataSaver[[]T], fileName string) error {
	var list []T
	for _, v := range d.data {
		list = append(list, v)
	}
	return saver.Save(&list, fileName)
}

func (d *DictionaryData[T]) Save(saver DataSaver[[]T]) error {
	return d.SaveTo(saver, d.fileName)
}