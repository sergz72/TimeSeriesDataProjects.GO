package entities

import "io"

type Category struct {
	Id   int
	Name string
}

func (c Category) GetId() int {
	return c.Id
}

func NewCategories(reader io.Reader) ([]Category, error) {
	return nil, nil
}
