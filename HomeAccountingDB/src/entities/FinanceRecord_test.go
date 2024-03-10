package entities

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"
)

func TestFinanceRecordBinary(t *testing.T) {
	var ops []FinanceOperation
	for i := 0; i < 100; i += 2 {
		amount := Decimal(i)
		s := strconv.Itoa(i)
		ops = append(ops, FinanceOperation{
			Date:            i,
			Amount:          &amount,
			Summa:           Decimal(i),
			SubcategoryId:   i,
			FinOpProperties: []FinOpProperty{{&i, &s, Date(i), Seca}},
			AccountId:       i,
		})
		ops = append(ops, FinanceOperation{
			Date:            i + 1,
			Amount:          nil,
			Summa:           Decimal(i + 1),
			SubcategoryId:   i + 1,
			FinOpProperties: []FinOpProperty{{nil, nil, Date(i + 1), Dist}},
			AccountId:       i + 1,
		})
	}
	r := NewFinanceRecord(ops)
	r.totals[1] = 2
	r.totals[2] = 3
	b := new(bytes.Buffer)
	err := r.Save(b)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := NewFinanceRecordFromBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(r, *r2) {
		t.Fatal("different objects")
	}
}
