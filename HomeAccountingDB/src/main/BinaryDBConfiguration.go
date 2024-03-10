package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"io"
	"strconv"
	"strings"
)

type binaryDatedSource struct {
	processor core.CryptoProcessor
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func (b *binaryDatedSource) GetFileDate(fileName string, _ string) (int, error) {
	return strconv.Atoi(fileNameWithoutExtension(fileName))
}

func (b *binaryDatedSource) Load(files []core.FileWithDate) (*entities.FinanceRecord, error) {
	return core.LoadBinaryP[entities.FinanceRecord](files[0].FileName, b.processor, entities.NewFinanceRecordFromBinary)
}

func (b *binaryDatedSource) getFileName(date int, dataFolderPath string) string {
	return dataFolderPath + "/" + strconv.Itoa(date) + ".bin"
}

func (b *binaryDatedSource) GetFiles(date int, dataFolderPath string) ([]core.FileWithDate, error) {
	return []core.FileWithDate{{
		FileName: b.getFileName(date, dataFolderPath),
		Date:     date,
	}}, nil
}

func (b *binaryDatedSource) Save(date int, data *entities.FinanceRecord, dataFolderPath string) error {
	return core.SaveBinary(b.getFileName(date, dataFolderPath), b.processor, data)
}

type binaryDBConfiguration struct {
	processor core.CryptoProcessor
}

func newBinaryDBConfiguration(processor core.CryptoProcessor) binaryDBConfiguration {
	return binaryDBConfiguration{processor: processor}
}

func (b binaryDBConfiguration) GetAccounts(fileName string) ([]entities.Account, error) {
	return core.LoadBinary[[]entities.Account](fileName, b.processor, func(reader io.Reader) ([]entities.Account, error) {
		return core.NewBinaryData[entities.Account](reader, entities.NewAccountFromBinary)
	})
}

func (b binaryDBConfiguration) GetCategories(fileName string) ([]entities.Category, error) {
	return core.LoadBinary[[]entities.Category](fileName, b.processor, func(reader io.Reader) ([]entities.Category, error) {
		return core.NewBinaryData[entities.Category](reader, entities.NewCategoryFromBinary)
	})
}

func (b binaryDBConfiguration) GetSubcategories(fileName string) ([]entities.Subcategory, error) {
	return core.LoadBinary[[]entities.Subcategory](fileName, b.processor, func(reader io.Reader) ([]entities.Subcategory, error) {
		return core.NewBinaryData[entities.Subcategory](reader, entities.NewSubcategoryFromBinary)
	})
}

func (b binaryDBConfiguration) GetMainDataSource() core.DatedSource[entities.FinanceRecord] {
	return &binaryDatedSource{}
}

func (b binaryDBConfiguration) GetAccountsSaver() core.DataSaver[[]entities.Account] {
	return core.BinarySaver[[]entities.Account]{}
}

func (b binaryDBConfiguration) GetCategoriesSaver() core.DataSaver[[]entities.Category] {
	return core.BinarySaver[[]entities.Category]{}
}

func (b binaryDBConfiguration) GetSubcategoriesSaver() core.DataSaver[[]entities.Subcategory] {
	return core.BinarySaver[[]entities.Subcategory]{}
}
