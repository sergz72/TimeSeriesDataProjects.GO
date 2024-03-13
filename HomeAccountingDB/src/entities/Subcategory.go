package entities

import (
	"TimeSeriesData/core"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

type SubcategoryCode int
type SubcategoryOperationCode int

const (
	Comb SubcategoryCode = 0
	Comc SubcategoryCode = iota
	Fuel SubcategoryCode = iota
	Prcn SubcategoryCode = iota
	Incc SubcategoryCode = iota
	Expc SubcategoryCode = iota
	Exch SubcategoryCode = iota
	Trfr SubcategoryCode = iota
	None SubcategoryCode = iota
)

const (
	Incm SubcategoryOperationCode = 0
	Expn SubcategoryOperationCode = iota
	Spcl SubcategoryOperationCode = iota
)

func (n *SubcategoryCode) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*n = None
		return nil
	}
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	switch v {
	case "COMB":
		*n = Comb
	case "COMC":
		*n = Comc
	case "FUEL":
		*n = Fuel
	case "PRCN":
		*n = Prcn
	case "INCC":
		*n = Incc
	case "EXPC":
		*n = Expc
	case "EXCH":
		*n = Exch
	case "TRFR":
		*n = Trfr
	default:
		return errors.New("unknown subcategory code")
	}
	return nil
}

func (n *SubcategoryOperationCode) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	switch v {
	case "INCM":
		*n = Incm
	case "EXPN":
		*n = Expn
	case "SPCL":
		*n = Spcl
	default:
		return errors.New("unknown subcategory operation code")
	}
	return nil
}

type Subcategory struct {
	Id              int
	Code            SubcategoryCode
	Name            string
	OperationCodeId SubcategoryOperationCode
	CategoryId      int
}

func (s Subcategory) GetId() int {
	return s.Id
}

func (s Subcategory) Save(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint32(s.Id))
	if err != nil {
		return err
	}
	err = core.WriteStringToBinary(writer, s.Name)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint32(s.CategoryId))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint8(s.Code))
	if err != nil {
		return err
	}
	return binary.Write(writer, binary.LittleEndian, uint8(s.OperationCodeId))
}

func NewSubcategoryFromBinary(reader io.Reader) (Subcategory, error) {
	var id uint32
	err := binary.Read(reader, binary.LittleEndian, &id)
	if err != nil {
		return Subcategory{}, nil
	}
	var name string
	name, err = core.ReadStringFromBinary(reader)
	if err != nil {
		return Subcategory{}, nil
	}
	var categoryId uint32
	err = binary.Read(reader, binary.LittleEndian, &categoryId)
	if err != nil {
		return Subcategory{}, nil
	}
	var code uint8
	err = binary.Read(reader, binary.LittleEndian, &code)
	if err != nil {
		return Subcategory{}, nil
	}
	var operationCode uint8
	err = binary.Read(reader, binary.LittleEndian, &operationCode)
	if err != nil {
		return Subcategory{}, nil
	}
	return Subcategory{Id: int(id), Name: name, CategoryId: int(categoryId), Code: SubcategoryCode(code),
		OperationCodeId: SubcategoryOperationCode(operationCode)}, nil
}

func SaveSubcategoryByIndex(index int, value []Subcategory, writer io.Writer) error {
	return value[index].Save(writer)
}
