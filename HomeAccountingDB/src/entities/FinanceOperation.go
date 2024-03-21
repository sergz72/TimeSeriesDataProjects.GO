package entities

import (
	"TimeSeriesData/core"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
)

type FinOpPropertyCode int

const (
	Amou FinOpPropertyCode = 0
	Dist FinOpPropertyCode = iota
	Netw FinOpPropertyCode = iota
	Ppto FinOpPropertyCode = iota
	Seca FinOpPropertyCode = iota
	Typ  FinOpPropertyCode = iota
)

func FinOpPropertyCodeFromString(v string) (FinOpPropertyCode, error) {
	switch v {
	case "AMOU":
		return Amou, nil
	case "DIST":
		return Dist, nil
	case "NETW":
		return Netw, nil
	case "PPTO":
		return Ppto, nil
	case "SECA":
		return Seca, nil
	case "TYPE":
		return Typ, nil
	default:
		return Amou, errors.New("unknown fin_op_property code")
	}
}

func (n *FinOpPropertyCode) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	*n, err = FinOpPropertyCodeFromString(v)
	return err
}

type FinOpProperty struct {
	NumericValue *int
	StringValue  *string
	DateValue    Date
	PropertyCode FinOpPropertyCode
}

type FinanceOperation struct {
	Date            int `json:"-"`
	Amount          *Decimal
	Summa           Decimal
	SubcategoryId   int
	FinOpProperties []FinOpProperty `json:"finOpProperies"`
	AccountId       int
}

type FinanceChange struct {
	StartBalance     int
	SummaIncome      int
	SummaExpenditure int
}

func (c *FinanceChange) GetEndSumma() int {
	return c.StartBalance + c.SummaIncome - c.SummaExpenditure
}

func (c *FinanceChange) Income(summa int) {
	c.SummaIncome += summa
}

func (c *FinanceChange) Expenditure(summa int) {
	c.SummaExpenditure += summa
}

func (op *FinanceOperation) UpdateChanges(changes map[int]*FinanceChange, accounts core.DictionaryData[Account],
	subcategories core.DictionaryData[Subcategory]) error {
	subcategory, err := subcategories.Get(op.SubcategoryId)
	if err != nil {
		return err
	}
	change, ok := changes[op.AccountId]
	if !ok {
		change = &FinanceChange{}
		changes[op.AccountId] = change
	}
	switch subcategory.OperationCodeId {
	case Incm:
		change.Income(int(op.Summa))
	case Expn:
		change.Expenditure(int(op.Summa))
	case Spcl:
		switch subcategory.Code {
		case Incc: // Пополнение карточного счета наличными
			return op.handleINCC(change, changes, accounts)
		case Expc: // Снятие наличных в банкомате
			return op.handleEXPC(change, changes, accounts)
		case Exch: // Обмен валюты
			return op.handleEXCH(change, changes)
		case Trfr: // Перевод средств между платежными картами
			return op.handleTRFR(change, changes, int(op.Summa))
		default:
			return fmt.Errorf("invalid subcategory code: %v", subcategory.Code)
		}
	}
	return nil
}

func getCashAccount(id int, accounts core.DictionaryData[Account]) (int, error) {
	acc, err := accounts.Get(id)
	if err != nil {
		return 0, err
	}
	if acc.CashAccount <= 0 {
		return 0, errors.New("no cash account found")
	}
	return int(acc.CashAccount), nil
}

// Пополнение карточного счета наличными
func (op *FinanceOperation) handleINCC(current *FinanceChange, changes map[int]*FinanceChange, accounts core.DictionaryData[Account]) error {
	current.Income(int(op.Summa))
	// cash account for corresponding currency
	accountId, err := getCashAccount(op.AccountId, accounts)
	if err != nil {
		return err
	}
	change2, ok := changes[accountId]
	if !ok {
		change2 = &FinanceChange{}
		changes[accountId] = change2
	}
	change2.Expenditure(int(op.Summa))
	return nil
}

// Снятие наличных в банкомате
func (op *FinanceOperation) handleEXPC(current *FinanceChange, changes map[int]*FinanceChange, accounts core.DictionaryData[Account]) error {
	current.Expenditure(int(op.Summa))
	// cash account for corresponding currency
	accountId, err := getCashAccount(op.AccountId, accounts)
	if err != nil {
		return err
	}
	change2, ok := changes[accountId]
	if !ok {
		change2 = &FinanceChange{}
		changes[accountId] = change2
	}
	change2.Income(int(op.Summa))
	return nil
}

// Обмен валюты
func (op *FinanceOperation) handleEXCH(current *FinanceChange, changes map[int]*FinanceChange) error {
	if op.Amount == nil {
		return nil
	}
	return op.handleTRFR(current, changes, int(*op.Amount)/10)
}

// Перевод средств между платежными картами
func (op *FinanceOperation) handleTRFR(current *FinanceChange, changes map[int]*FinanceChange, summa int) error {
	current.Expenditure(summa)
	if op.FinOpProperties != nil {
		for _, property := range op.FinOpProperties {
			if property.PropertyCode == Seca {
				if property.NumericValue == nil {
					return nil
				}
				accountId := *property.NumericValue
				change2, ok := changes[accountId]
				if !ok {
					change2 = &FinanceChange{}
					changes[accountId] = change2
				}
				change2.Income(int(op.Summa))
				return nil
			}
		}
	}
	return nil
}

func (prop *FinOpProperty) SaveToBinary(writer io.Writer) error {
	var err error
	if prop.NumericValue != nil {
		err = binary.Write(writer, binary.LittleEndian, int64(*prop.NumericValue))
	} else {
		err = binary.Write(writer, binary.LittleEndian, int64(math.MaxInt64))
	}
	if err != nil {
		return err
	}
	if prop.StringValue != nil {
		err = core.WriteStringToBinary(writer, *prop.StringValue)
	} else {
		err = core.WriteStringToBinary(writer, "")
	}
	if err != nil {
		return err
	}
	v := uint32(prop.DateValue)
	err = binary.Write(writer, binary.LittleEndian, v)
	if err != nil {
		return err
	}
	c := uint8(prop.PropertyCode)
	return binary.Write(writer, binary.LittleEndian, c)
}

func (op *FinanceOperation) SaveToBinary(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint32(op.Date))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint32(op.AccountId))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint32(op.SubcategoryId))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, int64(op.Summa))
	if err != nil {
		return err
	}
	var a int64 = math.MaxInt64
	if op.Amount != nil {
		a = int64(*op.Amount)
	}
	err = binary.Write(writer, binary.LittleEndian, a)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint32(len(op.FinOpProperties)))
	if err != nil {
		return err
	}
	for _, prop := range op.FinOpProperties {
		err = prop.SaveToBinary(writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewFinOpPropertyFromBinary(reader io.Reader) (FinOpProperty, error) {
	var prop FinOpProperty
	var v64 int64
	err := binary.Read(reader, binary.LittleEndian, &v64)
	if v64 != math.MaxInt64 {
		v := int(v64)
		prop.NumericValue = &v
	}
	var s string
	s, err = core.ReadStringFromBinary(reader)
	if len(s) != 0 {
		prop.StringValue = &s
	}
	var d uint32
	err = binary.Read(reader, binary.LittleEndian, &d)
	if err != nil {
		return prop, err
	}
	prop.DateValue = Date(d)
	var c uint8
	err = binary.Read(reader, binary.LittleEndian, &c)
	if err != nil {
		return prop, err
	}
	prop.PropertyCode = FinOpPropertyCode(c)
	return prop, nil
}

func NewFinanceOperationFromBinary(reader io.Reader) (FinanceOperation, error) {
	var op FinanceOperation
	var v uint32
	var v64 int64
	err := binary.Read(reader, binary.LittleEndian, &v)
	if err != nil {
		return op, err
	}
	op.Date = int(v)
	err = binary.Read(reader, binary.LittleEndian, &v)
	if err != nil {
		return op, err
	}
	op.AccountId = int(v)
	err = binary.Read(reader, binary.LittleEndian, &v)
	if err != nil {
		return op, err
	}
	op.SubcategoryId = int(v)
	err = binary.Read(reader, binary.LittleEndian, &v64)
	if err != nil {
		return op, err
	}
	op.Summa = Decimal(v64)
	err = binary.Read(reader, binary.LittleEndian, &v64)
	if err != nil {
		return op, err
	}
	if v64 != math.MaxInt64 {
		var d = Decimal(v64)
		op.Amount = &d
	}
	var ll uint32
	err = binary.Read(reader, binary.LittleEndian, &ll)
	if err != nil {
		return op, err
	}
	for ll > 0 {
		var prop FinOpProperty
		prop, err = NewFinOpPropertyFromBinary(reader)
		op.FinOpProperties = append(op.FinOpProperties, prop)
		ll--
	}
	return op, nil
}
