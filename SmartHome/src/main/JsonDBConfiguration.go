package main

import (
	"SmartHome/src/entities"
	"TimeSeriesData/core"
	"errors"
	"strconv"
)

type jsonDatedSource struct{}

func (s jsonDatedSource) GetFileDate(_ string, folderName string) (int, error) {
	return strconv.Atoi(folderName)
}

func (s jsonDatedSource) Load(files []core.FileWithDate) (*entities.SensorData, error) {
	var data map[int][]entities.SensorDataItem
	for _, f := range files {
		sensorId, err := strconv.Atoi(fileNameWithoutExtension(f.FileName))
		if err != nil {
			return nil, err
		}
		items, err := core.LoadJson[[]entities.SensorDataItem](f.FileName)
		if err != nil {
			return nil, err
		}
		sensorData, ok := data[sensorId]
		if ok {
			data[sensorId] = append(sensorData, items...)
		} else {
			data[sensorId] = items
		}
	}
	return entities.NewSensorData(data), nil
}

func (s jsonDatedSource) GetFiles(date int, dataFolderPath string) ([]core.FileWithDate, error) {
	return nil, errors.New("not implemented")
}

func (s jsonDatedSource) Save(date int, data *entities.SensorData, dataFolderPath string) error {
	return errors.New("not implemented")
}

type jsonDBConfiguration struct{}

func (c jsonDBConfiguration) GetSensors(fileName string) ([]entities.Sensor, error) {
	return core.LoadJson[[]entities.Sensor](fileName + ".json")
}

func (c jsonDBConfiguration) GetLocations(fileName string) ([]entities.Location, error) {
	return core.LoadJson[[]entities.Location](fileName + ".json")
}

func (c jsonDBConfiguration) GetMainDataSource() core.DatedSource[entities.SensorData] {
	return jsonDatedSource{}
}
