package entities

import (
	"TimeSeriesData/core"
	"encoding/binary"
	"encoding/json"
	"io"
)

type Int int

func (n *Int) UnmarshalJSON(b []byte) error {
	var v bool
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	if v {
		*n = -1
	} else {
		*n = 0
	}
	return nil
}

type Account struct {
	Id          int
	Name        string
	CashAccount Int `json:"isCash"`
	ActiveTo    Date
	Currency    string `json:"valutaCode"`
}

func (a Account) GetId() int {
	return a.Id
}

func (a *Account) GetName() string {
	return a.Name
}

func (a Account) Save(writer io.Writer) error {
	err := binary.Write(writer, binary.BigEndian, uint32(a.Id))
	if err != nil {
		return err
	}
	err = core.WriteStringToBinary(writer, a.Name)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.BigEndian, int32(a.CashAccount))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.BigEndian, uint32(a.ActiveTo))
	if err != nil {
		return err
	}
	return core.WriteStringToBinary(writer, a.Currency)
}

func NewAccountFromBinary(reader io.Reader) (Account, error) {
	return Account{}, nil
}
