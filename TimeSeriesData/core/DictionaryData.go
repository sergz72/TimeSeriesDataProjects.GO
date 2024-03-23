package core

import (
	"errors"
	"io"
	"os"
)

type DataSource[T any] interface {
	Load(fileName string) (T, error)
}

type DataSaver interface {
	Save(data any, saveIndex func(int, any, io.Writer) error) error
	GetBytes() []byte
	GetFileExtension() string
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

func (d *DictionaryData[T]) SaveTo(saver DataSaver, saveIndex func(int, any, io.Writer) error) error {
	var list []T
	for _, v := range d.data {
		list = append(list, v)
	}
	return saver.Save(list, saveIndex)
}

func (d *DictionaryData[T]) SaveToFile(saver DataSaver, fileName string, saveIndex func(int, any, io.Writer) error) error {
	err := d.SaveTo(saver, saveIndex)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName+saver.GetFileExtension(), saver.GetBytes(), 0644)
}

func (d *DictionaryData[T]) Save(saver DataSaver, fileName string, saveIndex func(int, any, io.Writer) error) error {
	return d.SaveToFile(saver, fileName, saveIndex)
}
