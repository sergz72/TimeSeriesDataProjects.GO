package entities

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSubcategoryBinary(t *testing.T) {
	s := Subcategory{
		Id:              1,
		Code:            Exch,
		Name:            "subcategory1",
		OperationCodeId: Expn,
		CategoryId:      3,
	}
	buffer := new(bytes.Buffer)
	err := s.Save(buffer)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := NewSubcategoryFromBinary(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(s, s2) {
		t.Fatal("different subcategories")
	}
	if buffer.Len() != 0 {
		t.Fatal("non zero length")
	}
}
