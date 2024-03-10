package entities

import (
	"TimeSeriesData/core"
	"encoding/binary"
	"io"
)

type Category struct {
	Id   int
	Name string
}

func (c Category) GetId() int {
	return c.Id
}

func (c Category) Save(writer io.Writer) error {
	err := binary.Write(writer, binary.BigEndian, uint32(c.Id))
	if err != nil {
		return err
	}
	return core.WriteStringToBinary(writer, c.Name)
}

func NewCategoryFromBinary(reader io.Reader) (Category, error) {
	return Category{}, nil
}
