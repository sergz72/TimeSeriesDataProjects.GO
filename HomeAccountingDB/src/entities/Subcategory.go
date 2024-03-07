package entities

import (
	"encoding/json"
	"errors"
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
