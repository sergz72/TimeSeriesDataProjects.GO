package entities

import (
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

func NewAccounts(reader io.Reader) ([]Account, error) {
	return nil, nil
}

type Accounts []Account

func (a Accounts) Save(writer io.Writer) error {
	return nil
}
