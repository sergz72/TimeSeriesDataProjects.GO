package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"encoding/binary"
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
	date, err := strconv.Atoi(fileNameWithoutExtension(fileName))
	return date * 100, err
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

type hintsItem struct {
	code  entities.FinOpPropertyCode
	hints map[string]bool
}

func newHintsItem(reader io.Reader) (hintsItem, error) {
	var code uint8
	err := binary.Read(reader, binary.LittleEndian, &code)
	if err != nil {
		return hintsItem{}, err
	}
	var l uint8
	err = binary.Read(reader, binary.LittleEndian, &l)
	if err != nil {
		return hintsItem{}, err
	}
	hints := make(map[string]bool)
	for l > 0 {
		var hint string
		hint, err = core.ReadStringFromBinary(reader)
		if err != nil {
			return hintsItem{}, err
		}
		hints[hint] = true
		l--
	}
	return hintsItem{code: entities.FinOpPropertyCode(code), hints: hints}, nil
}

func (b binaryDBConfiguration) getHintsFromData(data []byte) (dbHints, error) {
	items, err := core.LoadBinaryData[[]hintsItem](data, b.processor, func(reader io.Reader) ([]hintsItem, error) {
		return core.LoadBinaryArray[hintsItem](reader, newHintsItem)
	})
	if err != nil {
		return nil, err
	}
	result := make(map[entities.FinOpPropertyCode]map[string]bool)
	for _, item := range items {
		result[item.code] = item.hints
	}
	return result, nil
}

func (b binaryDBConfiguration) GetHints(fileName string) (dbHints, error) {
	items, err := core.LoadBinary[[]hintsItem](fileName+".bin", b.processor, func(reader io.Reader) ([]hintsItem, error) {
		return core.LoadBinaryArray[hintsItem](reader, newHintsItem)
	})
	if err != nil {
		return nil, err
	}
	result := make(map[entities.FinOpPropertyCode]map[string]bool)
	for _, item := range items {
		result[item.code] = item.hints
	}
	return result, nil
}

func (h dbHints) Save(writer io.Writer) error {
	l := uint16(len(h))
	err := binary.Write(writer, binary.LittleEndian, l)
	if err != nil {
		return err
	}
	for k, v := range h {
		err = binary.Write(writer, binary.LittleEndian, uint8(k))
		if err != nil {
			return err
		}
		err = binary.Write(writer, binary.LittleEndian, uint8(len(v)))
		if err != nil {
			return err
		}
		for hint := range v {
			err = core.WriteStringToBinary(writer, hint)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b binaryDBConfiguration) GetHintsSaver() core.DataSaver[dbHints] {
	return core.NewBinarySaver[dbHints](b.processor)
}

func newBinaryDBConfiguration(processor core.CryptoProcessor) binaryDBConfiguration {
	return binaryDBConfiguration{processor: processor}
}

func (b binaryDBConfiguration) GetAccounts(fileName string) ([]entities.Account, error) {
	return core.LoadBinary[[]entities.Account](fileName+".bin", b.processor, func(reader io.Reader) ([]entities.Account, error) {
		return core.LoadBinaryArray[entities.Account](reader, entities.NewAccountFromBinary)
	})
}

func (b binaryDBConfiguration) GetCategories(fileName string) ([]entities.Category, error) {
	return core.LoadBinary[[]entities.Category](fileName+".bin", b.processor, func(reader io.Reader) ([]entities.Category, error) {
		return core.LoadBinaryArray[entities.Category](reader, entities.NewCategoryFromBinary)
	})
}

func (b binaryDBConfiguration) GetSubcategories(fileName, mapFileName string) ([]entities.Subcategory, error) {
	return core.LoadBinary[[]entities.Subcategory](fileName+".bin", b.processor, func(reader io.Reader) ([]entities.Subcategory, error) {
		return core.LoadBinaryArray[entities.Subcategory](reader, entities.NewSubcategoryFromBinary)
	})
}

func (b binaryDBConfiguration) GetMainDataSource() core.DatedSource[entities.FinanceRecord] {
	return &binaryDatedSource{b.processor}
}

func (b binaryDBConfiguration) GetAccountsSaver() core.DataSaver[[]entities.Account] {
	return core.NewBinarySaver[[]entities.Account](b.processor)
}

func (b binaryDBConfiguration) GetCategoriesSaver() core.DataSaver[[]entities.Category] {
	return core.NewBinarySaver[[]entities.Category](b.processor)
}

func (b binaryDBConfiguration) GetSubcategoriesSaver() core.DataSaver[[]entities.Subcategory] {
	return core.NewBinarySaver[[]entities.Subcategory](b.processor)
}
