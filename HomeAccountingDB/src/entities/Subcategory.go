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

func SubcategoryCodeFromString(v string) (SubcategoryCode, error) {
	switch v {
	case "COMB":
		return Comb, nil
	case "COMC":
		return Comc, nil
	case "FUEL":
		return Fuel, nil
	case "PRCN":
		return Prcn, nil
	case "INCC":
		return Incc, nil
	case "EXPC":
		return Expc, nil
	case "EXCH":
		return Exch, nil
	case "TRFR":
		return Trfr, nil
	default:
		return None, errors.New("unknown subcategory code")
	}
}

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
	*n, err = SubcategoryCodeFromString(v)
	return err
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
	Id                 int
	Code               SubcategoryCode
	Name               string
	OperationCodeId    SubcategoryOperationCode
	CategoryId         int
	RequiredProperties []FinOpPropertyCode
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
	err = binary.Write(writer, binary.LittleEndian, uint8(s.OperationCodeId))
	if err != nil {
		return err
	}
	l := uint8(len(s.RequiredProperties))
	err = binary.Write(writer, binary.LittleEndian, l)
	if err != nil {
		return err
	}
	for _, prop := range s.RequiredProperties {
		err = binary.Write(writer, binary.LittleEndian, uint8(prop))
		if err != nil {
			return err
		}
	}
	return nil
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
	var l uint8
	err = binary.Read(reader, binary.LittleEndian, &l)
	if err != nil {
		return Subcategory{}, nil
	}
	var requiredProperties []FinOpPropertyCode
	for l > 0 {
		var prop uint8
		err = binary.Read(reader, binary.LittleEndian, &prop)
		if err != nil {
			return Subcategory{}, nil
		}
		requiredProperties = append(requiredProperties, FinOpPropertyCode(prop))
		l--
	}
	return Subcategory{Id: int(id), Name: name, CategoryId: int(categoryId), Code: SubcategoryCode(code),
		OperationCodeId: SubcategoryOperationCode(operationCode), RequiredProperties: requiredProperties}, nil
}

func SaveSubcategoryByIndex(index int, value any, writer io.Writer) error {
	v := value.([]Subcategory)
	return v[index].Save(writer)
}
