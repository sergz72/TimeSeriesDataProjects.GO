package main

import (
	"fmt"
	"testing"
)

func TestFrom(t *testing.T) {
	converter := newDateConverter(1970, 1, 2)
	v := converter.fromDate(19710302)
	shouldBe := 365 + 31 + 28 + 1
	if v != shouldBe {
		t.Fatal(fmt.Printf("fromDate 1970 1 19710302 error %v %v", v, shouldBe))
	}
}
