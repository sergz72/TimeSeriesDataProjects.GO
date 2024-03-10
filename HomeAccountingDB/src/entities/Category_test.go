package entities

import (
	"bytes"
	"reflect"
	"testing"
)

func TestCategoryBinary(t *testing.T) {
	c := Category{
		Id:   1,
		Name: "category1",
	}
	buffer := new(bytes.Buffer)
	err := c.Save(buffer)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := NewCategoryFromBinary(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, c2) {
		t.Fatal("different categories")
	}
	if buffer.Len() != 0 {
		t.Fatal("non zero length")
	}
}
