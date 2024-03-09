package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
)

type binaryDatedSource struct{}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func (b binaryDatedSource) GetFileDate(fileName string, folderName string) (int, error) {
	return strconv.Atoi(fileNameWithoutExtension(fileName))
}

func (b binaryDatedSource) Load(files []core.FileWithDate) (*entities.FinanceRecord, error) {
	file, err := os.Open(files[0].FileName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	binary.Read(file, binary.BigEndian, 13)
	buffer := bytes.NewBuffer(data)
	return entities.NewFinanceRecord(file)
}

func (b binaryDatedSource) GetFiles(date int) ([]core.FileWithDate, error) {
	return []core.FileWithDate{{
		FileName: strconv.Itoa(date) + ".bin",
		Date:     date,
	}}, nil
}

func (b binaryDatedSource) Save(date int, data *entities.FinanceRecord) error {
	//TODO implement me
	panic("implement me")
}

type binaryDBConfiguration struct{}

func (b binaryDBConfiguration) GetAccounts(fileName string) ([]entities.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (b binaryDBConfiguration) GetCategories(fileName string) ([]entities.Category, error) {
	//TODO implement me
	panic("implement me")
}

func (b binaryDBConfiguration) GetSubcategories(fileName string) ([]entities.Subcategory, error) {
	//TODO implement me
	panic("implement me")
}

func (b binaryDBConfiguration) GetMainDataSource() core.DatedSource[entities.FinanceRecord] {
	return binaryDatedSource{}
}

func (b binaryDBConfiguration) GetAccountsSaver() core.DataSaver[[]entities.Account] {
	//TODO implement me
	panic("implement me")
}

func (b binaryDBConfiguration) GetCategoriesSaver() core.DataSaver[[]entities.Category] {
	//TODO implement me
	panic("implement me")
}

func (b binaryDBConfiguration) GetSubcategoriesSaver() core.DataSaver[[]entities.Subcategory] {
	//TODO implement me
	panic("implement me")
}
