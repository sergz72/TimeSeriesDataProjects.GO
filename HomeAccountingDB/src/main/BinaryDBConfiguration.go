package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"strconv"
	"strings"
)

type binaryDatedSource struct {
	dataFolderPath string
	processor      core.CryptoProcessor
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func (b *binaryDatedSource) GetFileDate(fileName string, folderName string) (int, error) {
	return strconv.Atoi(fileNameWithoutExtension(fileName))
}

func (b *binaryDatedSource) Load(files []core.FileWithDate) (*entities.FinanceRecord, error) {
	return core.LoadBinaryP[entities.FinanceRecord](files[0].FileName, b.processor, entities.NewFinanceRecordFromBinary)
}

func (b *binaryDatedSource) getFileName(date int) string {
	return b.dataFolderPath + "/" + strconv.Itoa(date) + ".bin"
}

func (b *binaryDatedSource) GetFiles(date int) ([]core.FileWithDate, error) {
	return []core.FileWithDate{{
		FileName: b.getFileName(date),
		Date:     date,
	}}, nil
}

func (b *binaryDatedSource) Save(date int, data *entities.FinanceRecord) error {
	return core.SaveBinary(b.getFileName(date), b.processor, data)
}

type binaryDBConfiguration struct {
	processor core.CryptoProcessor
}

func (b binaryDBConfiguration) GetAccounts(fileName string) ([]entities.Account, error) {
	return core.LoadBinary[[]entities.Account](fileName, b.processor, core.CreateFromBinary[[]entities.Account])
}

func (b binaryDBConfiguration) GetCategories(fileName string) ([]entities.Category, error) {
	return core.LoadBinary[[]entities.Category](fileName, b.processor, core.CreateFromBinary[[]entities.Category])
}

func (b binaryDBConfiguration) GetSubcategories(fileName string) ([]entities.Subcategory, error) {
	return core.LoadBinary[[]entities.Subcategory](fileName, b.processor, core.CreateFromBinary[[]entities.Subcategory])
}

func (b binaryDBConfiguration) GetMainDataSource() core.DatedSource[entities.FinanceRecord] {
	return &binaryDatedSource{}
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
