package main

import (
	"SmartHome/src/entities"
	"TimeSeriesData/core"
	"fmt"
)

type settings struct {
	MinYear                  int
	MinMonth                 int
	TimeSeriesDataCapacity   int
	MaxActiveTimeSeriesItems int
	YearsToCreate            int
	DataFolderPath           string
	ServerPort               int
}

type dBConfiguration interface {
	GetSensors(fileName string) ([]entities.Sensor, error)
	GetLocations(fileName string) ([]entities.Location, error)
	GetMainDataSource() core.DatedSource[entities.SensorData]
}

type dB struct {
	converter      dateConverter
	dataFolderPath string
	configuration  dBConfiguration
	sensors        core.DictionaryData[entities.Sensor]
	locations      core.DictionaryData[entities.Location]
	data           core.TimeSeriesData[entities.SensorData]
}

func getSensorsFileName(dataFolderPath string) string {
	return dataFolderPath + "/sensors"
}

func getLocationsFileName(dataFolderPath string) string {
	return dataFolderPath + "/locations"
}

func getMainDataFolderPath(dataFolderPath string) string {
	return dataFolderPath + "/dates"
}

func loadDicts(s settings, configuration dBConfiguration) (core.DictionaryData[entities.Sensor],
	core.DictionaryData[entities.Location], error) {
	path := getSensorsFileName(s.DataFolderPath)
	slist, err := configuration.GetSensors(path)
	if err != nil {
		return core.DictionaryData[entities.Sensor]{}, core.DictionaryData[entities.Location]{}, err
	}
	sensors := core.NewDictionaryData[entities.Sensor](path, "sensor", slist)
	path = getLocationsFileName(s.DataFolderPath)
	llist, err := configuration.GetLocations(path)
	if err != nil {
		return core.DictionaryData[entities.Sensor]{}, core.DictionaryData[entities.Location]{}, err
	}
	locations := core.NewDictionaryData[entities.Location](path, "location", llist)
	return sensors, locations, nil
}

func loadDB(s settings, configuration dBConfiguration) (*dB, error) {
	sensors, locations, err := loadDicts(s, configuration)
	if err != nil {
		return nil, err
	}
	converter := newDateConverter(s.MinYear, s.MinMonth, s.YearsToCreate)
	data, err := core.LoadTimeSeriesData[entities.SensorData](getMainDataFolderPath(s.DataFolderPath),
		configuration.GetMainDataSource(), s.TimeSeriesDataCapacity, func(date int) int {
			return converter.fromDate(date)
		}, func(idx int) int {
			return converter.toDate(idx)
		}, s.MaxActiveTimeSeriesItems)
	if err != nil {
		return nil, err
	}
	return &dB{dataFolderPath: s.DataFolderPath, configuration: configuration, sensors: sensors, locations: locations,
		data: data}, nil
}

func initDB(s settings, configuration dBConfiguration) (*dB, error) {
	sensors, locations, err := loadDicts(s, configuration)
	if err != nil {
		return nil, err
	}
	converter := newDateConverter(s.MinYear, s.MinMonth, s.YearsToCreate)
	data, err := core.InitTimeSeriesData[entities.SensorData](getMainDataFolderPath(s.DataFolderPath),
		configuration.GetMainDataSource(), s.TimeSeriesDataCapacity, func(date int) int {
			return converter.fromDate(date)
		}, func(idx int) int {
			return converter.toDate(idx)
		}, s.MaxActiveTimeSeriesItems)
	if err != nil {
		return nil, err
	}
	return &dB{dataFolderPath: s.DataFolderPath, configuration: configuration, sensors: sensors, locations: locations,
		data: data}, nil
}

func (d *dB) printStats(date int) {
	v, err := d.data.GetExact(date)
	if err != nil {
		panic(err)
	}
	if v != nil {
		panic("not implemented yet")
	} else {
		fmt.Println("no data")
	}
}

func (d *dB) save() error {
	return d.saveTo(d.dataFolderPath, d.configuration)
}

func (d *dB) saveTo(dataFolderPath string, configuration dBConfiguration) error {
	return d.data.SaveAll(configuration.GetMainDataSource(), getMainDataFolderPath(dataFolderPath))
}
