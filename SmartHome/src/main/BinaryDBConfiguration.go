package main

import (
	"SmartHome/src/entities"
	"TimeSeriesData/core"
	"os"
	"strconv"
	"strings"
)

type binaryDatedSource struct {
}

func fileNameWithoutExtension(fileName string) string {
	start := strings.LastIndexByte(fileName, os.PathSeparator)
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[start+1 : pos]
	}
	return fileName
}

func (b *binaryDatedSource) GetFileDate(fileName string, _ string) (int, error) {
	date, err := strconv.Atoi(fileNameWithoutExtension(fileName))
	return date, err
}

func (b *binaryDatedSource) Load(files []core.FileWithDate) (*entities.SensorData, error) {
	return core.LoadBinaryP[entities.SensorData](files[0].FileName, nil, entities.NewSensorDataFromBinary)
}

func (b *binaryDatedSource) getFileName(date int, year string, dataFolderPath string) string {
	return dataFolderPath + "/" + year + "/" + strconv.Itoa(date) + ".bin"
}

func (b *binaryDatedSource) GetFiles(date int, dataFolderPath string) ([]core.FileWithDate, error) {
	year := strconv.Itoa(date / 10000)
	return []core.FileWithDate{{
		FileName: b.getFileName(date, year, dataFolderPath),
		Date:     date,
	}}, nil
}

func (b *binaryDatedSource) Save(date int, data *entities.SensorData, dataFolderPath string) error {
	year := strconv.Itoa(date / 10000)
	_ = os.Mkdir(dataFolderPath+"/"+year, 0700)
	return core.SaveBinary(b.getFileName(date, year, dataFolderPath), nil, data)
}

type binaryDBConfiguration struct {
}

func newBinaryDBConfiguration() binaryDBConfiguration {
	return binaryDBConfiguration{}
}

func (c binaryDBConfiguration) GetSensors(fileName string) ([]entities.Sensor, error) {
	return core.LoadJson[[]entities.Sensor](fileName + ".json")
}

func (c binaryDBConfiguration) GetLocations(fileName string) ([]entities.Location, error) {
	return core.LoadJson[[]entities.Location](fileName + ".json")
}

func (b binaryDBConfiguration) GetMainDataSource() core.DatedSource[entities.SensorData] {
	return &binaryDatedSource{}
}
