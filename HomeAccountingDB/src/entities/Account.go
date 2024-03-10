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
	var id uint32
	err := binary.Read(reader, binary.BigEndian, &id)
	if err != nil {
		return Account{}, nil
	}
	var name string
	name, err = core.ReadStringFromBinary(reader)
	if err != nil {
		return Account{}, nil
	}
	var cashAccount int32
	err = binary.Read(reader, binary.BigEndian, &cashAccount)
	if err != nil {
		return Account{}, nil
	}
	var activeTo uint32
	err = binary.Read(reader, binary.BigEndian, &activeTo)
	if err != nil {
		return Account{}, nil
	}
	var currency string
	currency, err = core.ReadStringFromBinary(reader)
	return Account{Id: int(id), Name: name, CashAccount: Int(cashAccount), ActiveTo: Date(activeTo), Currency: currency}, err
}
