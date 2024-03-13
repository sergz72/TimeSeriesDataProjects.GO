package entities

import (
	"TimeSeriesData/core"
	"encoding/binary"
	"fmt"
	"io"
)

type SensorDataItem struct {
	EventTime int
	Data      map[string]int
}

func (i *SensorDataItem) save(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint32(i.EventTime))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint8(len(i.Data)))
	if err != nil {
		return err
	}
	for dataType, value := range i.Data {
		err = core.WriteStringToBinary(writer, dataType)
		if err != nil {
			return err
		}
		err = binary.Write(writer, binary.LittleEndian, int32(value))
		if err != nil {
			return err
		}
	}
	return nil
}

func newSensorDataItemFromBinary(reader io.Reader) (SensorDataItem, error) {
	i := SensorDataItem{
		Data: make(map[string]int),
	}
	var eventTime uint32
	err := binary.Read(reader, binary.LittleEndian, &eventTime)
	if err != nil {
		return i, err
	}
	i.EventTime = int(eventTime)
	var l uint8
	err = binary.Read(reader, binary.LittleEndian, &l)
	if err != nil {
		return i, err
	}
	for l > 0 {
		var dataType string
		dataType, err = core.ReadStringFromBinary(reader)
		if err != nil {
			return i, err
		}
		var value int32
		err = binary.Read(reader, binary.LittleEndian, &value)
		if err != nil {
			return i, err
		}
		i.Data[dataType] = int(value)
		l--
	}
	return i, nil
}

type SensorDataStats struct {
	Min int
	Max int
	Avg int
	Sum int
	Cnt int
}

func newSensorDataStats(list []SensorDataItem) map[string]SensorDataStats {
	result := make(map[string]SensorDataStats)
	for _, item := range list {
		for dataType, value := range item.Data {
			dataTypeStats, ok := result[dataType]
			if ok {
				if dataTypeStats.Max < value {
					dataTypeStats.Max = value
				}
				if dataTypeStats.Min > value {
					dataTypeStats.Min = value
				}
				dataTypeStats.Cnt++
				dataTypeStats.Avg += value
				dataTypeStats.Sum += value
			} else {
				dataTypeStats = SensorDataStats{
					Min: value,
					Max: value,
					Avg: value,
					Sum: value,
					Cnt: 1,
				}
			}
			result[dataType] = dataTypeStats
		}
	}
	for dataType, stats := range result {
		stats.Avg /= stats.Cnt
		result[dataType] = stats
	}
	return result
}

type SensorData struct {
	data  map[int][]SensorDataItem
	stats map[int]map[string]SensorDataStats
}

func (s *SensorData) PrintStats() {
	for k, v := range s.data {
		fmt.Printf("%v %v\n", k, len(v))
	}
}

func aggregate(data map[int][]SensorDataItem) map[int]map[string]SensorDataStats {
	result := make(map[int]map[string]SensorDataStats)
	for sensorId, list := range data {
		result[sensorId] = newSensorDataStats(list)
	}
	return result
}

func (s *SensorData) Aggregate() {
	s.stats = aggregate(s.data)
}

func (s SensorData) Save(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint16(len(s.data)))
	if err != nil {
		return err
	}
	for sensorId, list := range s.data {
		err = binary.Write(writer, binary.LittleEndian, uint16(sensorId))
		if err != nil {
			return err
		}
		err = binary.Write(writer, binary.LittleEndian, uint16(len(list)))
		if err != nil {
			return err
		}
		for _, item := range list {
			err = item.save(writer)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewSensorData(data map[int][]SensorDataItem) *SensorData {
	return &SensorData{
		data:  data,
		stats: aggregate(data),
	}
}

func NewSensorDataFromBinary(reader io.Reader) (*SensorData, error) {
	var length uint16
	err := binary.Read(reader, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	data := make(map[int][]SensorDataItem)
	for length > 0 {
		var sensorId uint16
		err = binary.Read(reader, binary.LittleEndian, &sensorId)
		if err != nil {
			return nil, err
		}
		var listLength uint16
		err = binary.Read(reader, binary.LittleEndian, &listLength)
		if err != nil {
			return nil, err
		}
		var list []SensorDataItem
		for listLength > 0 {
			var item SensorDataItem
			item, err = newSensorDataItemFromBinary(reader)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
			listLength--
		}
		data[int(sensorId)] = list
		length--
	}
	return NewSensorData(data), nil
}
