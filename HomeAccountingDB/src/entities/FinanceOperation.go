package entities

import (
	"TimeSeriesData/core"
	"encoding/json"
	"errors"
	"fmt"
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

func (n *FinOpPropertyCode) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	switch v {
	case "AMOU":
		*n = Amou
	case "DIST":
		*n = Dist
	case "NETW":
		*n = Netw
	case "PPTO":
		*n = Ppto
	case "SECA":
		*n = Seca
	case "TYPE":
		*n = Typ
	default:
		return errors.New("unknown fin_op_property code")
	}
	return nil
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
