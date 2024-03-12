package entities

import "io"

type SensorDataItem struct {
	EventTime int
	Data      map[string]int
}

type SensorDataStats struct {
	Min int
	Max int
	Avg int
	Sum int
	Cnt int
}

type SensorData struct {
	data  map[int][]SensorDataItem
	stats map[int]SensorDataStats
}

func (s SensorData) Save(writer io.Writer) error {
	//TODO implement me
	panic("implement me")
}

func NewSensorData(data map[int][]SensorDataItem) *SensorData {
	//TODO implement me
	return &SensorData{
		data:  data,
		stats: nil,
	}
}

func NewSensorDataFromBinary(reader io.Reader) (*SensorData, error) {
	//TODO implement me
	return nil, nil
}
