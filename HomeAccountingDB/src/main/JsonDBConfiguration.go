package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"bytes"
	"errors"
	"strconv"
)

type jsonDatedSource struct{}

func (s jsonDatedSource) GetFileDate(_ string, folderName string) (int, error) {
	return strconv.Atoi(folderName)
}

func (s jsonDatedSource) Load(files []core.FileWithDate) (*entities.FinanceRecord, error) {
	var operations []entities.FinanceOperation
	for _, f := range files {
		ops, err := core.LoadJson[[]entities.FinanceOperation](f.FileName)
		if err != nil {
			return nil, err
		}
		for idx := range ops {
			ops[idx].Date = f.Date
		}
		operations = append(operations, ops...)
	}
	return entities.NewFinanceRecord(operations), nil
}

func (s jsonDatedSource) GetFiles(date int, dataFolderPath string) ([]core.FileWithDate, error) {
	return nil, errors.New("not implemented")
}

func (s jsonDatedSource) Save(date int, data *entities.FinanceRecord, dataFolderPath string) error {
	return errors.New("not implemented")
}

type jsonDBConfiguration struct{}

func (c jsonDBConfiguration) GetHints(fileName string) (dbHints, error) {
	return make(map[entities.FinOpPropertyCode]map[string]bool), nil
}

func (c jsonDBConfiguration) GetHintsSaver(buffer *bytes.Buffer) core.DataSaver[dbHints] {
	return nil
}

func (c jsonDBConfiguration) GetAccountsSaver(buffer *bytes.Buffer) core.DataSaver[[]entities.Account] {
	//TODO implement me
	panic("implement me")
}

func (c jsonDBConfiguration) GetCategoriesSaver(buffer *bytes.Buffer) core.DataSaver[[]entities.Category] {
	//TODO implement me
	panic("implement me")
}

func (c jsonDBConfiguration) GetSubcategoriesSaver(buffer *bytes.Buffer) core.DataSaver[[]entities.Subcategory] {
	//TODO implement me
	panic("implement me")
}

func (c jsonDBConfiguration) GetOpsAndChangesSaver() core.DataSaver[entities.OpsAndChanges] {
	//TODO implement me
	panic("implement me")
}

func (c jsonDBConfiguration) GetAccounts(fileName string) ([]entities.Account, error) {
	data, err := core.LoadJson[[]entities.Account](fileName + ".json")
	if err != nil {
		return nil, err
	}
	cashAccounts := make(map[string]int)
	for _, v := range data {
		if v.CashAccount == -1 {
			cashAccounts[v.Currency] = v.Id
		}
	}
	for idx, v := range data {
		if v.CashAccount == 0 {
			data[idx].CashAccount = entities.Int(cashAccounts[v.Currency])
		}
	}
	return data, nil
}

func (c jsonDBConfiguration) GetCategories(fileName string) ([]entities.Category, error) {
	return core.LoadJson[[]entities.Category](fileName + ".json")
}

type subcategoryMap struct {
	SubcategoryCode string
	PropertyCode    string
}

func (c jsonDBConfiguration) GetSubcategories(fileName, mapFileName string) ([]entities.Subcategory, error) {
	subcategories, err := core.LoadJson[[]entities.Subcategory](fileName + ".json")
	if err != nil {
		return nil, err
	}
	subcategoriesMap, err := core.LoadJson[[]subcategoryMap](mapFileName + ".json")
	if err != nil {
		return nil, err
	}
	for _, item := range subcategoriesMap {
		var code entities.SubcategoryCode
		code, err = entities.SubcategoryCodeFromString(item.SubcategoryCode)
		if err != nil {
			return nil, err
		}
		var prop entities.FinOpPropertyCode
		prop, err = entities.FinOpPropertyCodeFromString(item.PropertyCode)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(subcategories); i++ {
			s := &subcategories[i]
			if s.Code == code {
				s.RequiredProperties = append(s.RequiredProperties, prop)
			}
		}
	}
	return subcategories, nil
}

func (c jsonDBConfiguration) GetMainDataSource() core.DatedSource[entities.FinanceRecord] {
	return jsonDatedSource{}
}
