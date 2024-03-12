package main

import (
	"fmt"
	"testing"
	"time"
)

func TestFrom(t *testing.T) {
	converter := newDateConverter(1970, 1, 10)
	v := converter.fromDate(19710302)
	shouldBe := 365 + 31 + 28 + 1
	if v != shouldBe {
		t.Fatal(fmt.Printf("fromDate 1970 1 19710302 error %v %v", v, shouldBe))
	}
	d := converter.toDate(v)
	if d != 19710302 {
		t.Fatal(fmt.Printf("toDate 1970 1 19710302 error %v", d))
	}

	v = converter.fromDate(19790101)
	shouldBe = 3287
	if v != shouldBe {
		t.Fatal(fmt.Printf("fromDate 1970 1 19790101 error %v %v", v, shouldBe))
	}

	converter = newDateConverter(1970, 7, 2)
	v = converter.fromDate(19700715)
	shouldBe = 14
	if v != shouldBe {
		t.Fatal(fmt.Printf("fromDate 1970 7 19700715 error %v %v", v, shouldBe))
	}
	d = converter.toDate(v)
	if d != 19700715 {
		t.Fatal(fmt.Printf("toDate 1970 7 19700715 error %v", d))
	}

	v = converter.fromDate(19710302)
	shouldBe = 31 + 31 + 30 + 31 + 30 + 31 +
		31 + 28 + 1
	if v != shouldBe {
		t.Fatal(fmt.Printf("fromDate 1970 7 19710302 error %v %v", v, shouldBe))
	}
	d = converter.toDate(v)
	if d != 19710302 {
		t.Fatal(fmt.Printf("toDate 1970 7 19710302 error %v", d))
	}

	converter = newDateConverter(1972, 1, 10)
	v = converter.fromDate(19790101)
	shouldBe = 2557
	if v != shouldBe {
		t.Fatal(fmt.Printf("fromDate 1972 1 19790101 error %v %v", v, shouldBe))
	}
}

func TestFromTo(t *testing.T) {
	converter := newDateConverter(2010, 1, 10)
	date := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	day := 0
	for date.Year() < 2020 {
		d := date.Year()*10000 + int(date.Month()*100) + date.Day()
		idx := converter.fromDate(d)
		if idx != day {
			t.Fatal(fmt.Printf("%v %v %v", date, idx, day))
		}
		dd := converter.toDate(idx)
		if d != dd {
			t.Fatal(fmt.Printf("%v %v %v", date, d, dd))
		}
		date = date.AddDate(0, 0, 1)
		day++
	}
}
