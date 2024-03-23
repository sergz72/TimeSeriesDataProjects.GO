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
	err := binary.Write(writer, binary.LittleEndian, uint32(c.Id))
	if err != nil {
		return err
	}
	return core.WriteStringToBinary(writer, c.Name)
}

func NewCategoryFromBinary(reader io.Reader) (Category, error) {
	var id uint32
	err := binary.Read(reader, binary.LittleEndian, &id)
	if err != nil {
		return Category{}, nil
	}
	var name string
	name, err = core.ReadStringFromBinary(reader)
	return Category{int(id), name}, err
}

func SaveCategoryByIndex(index int, value any, writer io.Writer) error {
	v := value.([]Category)
	return v[index].Save(writer)
}
