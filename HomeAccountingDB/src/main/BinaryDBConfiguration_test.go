package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"reflect"
	"testing"
)

func TestHints(t *testing.T) {
	hints := make(dbHints)
	hints[entities.Typ] = map[string]bool{"Type1": true, "Type2": true}
	hints[entities.Netw] = map[string]bool{"Netw1": true, "Netw2": true}
	config := binaryDBConfiguration{}
	saver, _ := config.GetHintsSaver(nil).(*core.BinarySaver[dbHints])
	err := saver.Save(hints, nil)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := config.getHintsFromData(saver.GetBytes())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(hints, loaded) {
		t.Fatal("different data")
	}
}
