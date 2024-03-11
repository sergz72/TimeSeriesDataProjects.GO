package core

import (
	"errors"
	"io"
)

type DataSource[T any] interface {
	Load(fileName string) (T, error)
}

type DataSaver[T any] interface {
	Save(data T, fileName string, saveIndex func(int, T, io.Writer) error) error
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

func (d *DictionaryData[T]) SaveTo(saver DataSaver[[]T], fileName string, saveIndex func(int, []T, io.Writer) error) error {
	var list []T
	for _, v := range d.data {
		list = append(list, v)
	}
	return saver.Save(list, fileName, saveIndex)
}

func (d *DictionaryData[T]) Save(saver DataSaver[[]T], saveIndex func(int, []T, io.Writer) error) error {
	return d.SaveTo(saver, d.fileName, saveIndex)
}