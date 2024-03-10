package entities

import (
	"bytes"
	"reflect"
	"testing"
)

func TestAccountBinary(t *testing.T) {
	a := Account{
		Id:          1,
		Name:        "account1",
		CashAccount: 5,
		ActiveTo:    10,
		Currency:    "UAH",
	}
	buffer := new(bytes.Buffer)
	err := a.Save(buffer)
	if err != nil {
		t.Fatal(err)
	}
	a2, err := NewAccountFromBinary(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, a2) {
		t.Fatal("different accounts")
	}
	if buffer.Len() != 0 {
		t.Fatal("non zero length")
	}
}
