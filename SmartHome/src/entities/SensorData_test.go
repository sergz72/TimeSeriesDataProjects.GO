package entities

import (
	"bytes"
	"reflect"
	"testing"
)

func buildTestData() *SensorData {
	return NewSensorData(map[int][]SensorDataItem{
		1: {
			{EventTime: 0, Data: map[string]int{
				"humi": 5000,
				"temp": 2200,
			}},
			{EventTime: 10, Data: map[string]int{
				"humi": 4000,
				"temp": 2000,
			}},
			{EventTime: 20, Data: map[string]int{
				"humi": 3000,
				"temp": 1800,
			}},
		},
		2: {
			{EventTime: 0, Data: map[string]int{
				"pres": 150000,
				"temp": 1800,
			}},
			{EventTime: 20, Data: map[string]int{
				"pres": 160000,
				"temp": 1900,
			}},
			{EventTime: 40, Data: map[string]int{
				"pres": 170000,
				"temp": 2000,
			}},
		},
	})
}

func TestSensorData_Aggregate(t *testing.T) {
	data := buildTestData()
	if len(data.stats) != 2 {
		t.Fatal("len(data.stats)")
	}
	s1 := data.stats[1]
	if len(s1) != 2 {
		t.Fatal("len(s1)")
	}
	testStats(t, s1["humi"], 5000, 3000, 4000, 12000, 3)
	testStats(t, s1["temp"], 2200, 1800, 2000, 6000, 3)
	s2 := data.stats[2]
	if len(s2) != 2 {
		t.Fatal("len(s2)")
	}
	testStats(t, s2["pres"], 170000, 150000, 160000, 480000, 3)
	testStats(t, s2["temp"], 2000, 1800, 1900, 5700, 3)
}

func testStats(t *testing.T, s SensorDataStats, max, min, avg, sum, cnt int) {
	if s.Max != max {
		t.Fatal("Max")
	}
	if s.Min != min {
		t.Fatal("Min")
	}
	if s.Avg != avg {
		t.Fatal("Avg")
	}
	if s.Sum != sum {
		t.Fatal("Sum")
	}
	if s.Cnt != cnt {
		t.Fatal("Cnt")
	}
}

func TestNewSensorDataBinary(t *testing.T) {
	data := buildTestData()
	buffer := new(bytes.Buffer)
	err := data.Save(buffer)
	if err != nil {
		t.Fatal(err)
	}
	data2, err := NewSensorDataFromBinary(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if buffer.Len() != 0 {
		t.Fatal("buffer.Len() should be 0")
	}
	if !reflect.DeepEqual(data, data2) {
		t.Fatal("different data")
	}
}
