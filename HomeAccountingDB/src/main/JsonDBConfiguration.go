package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
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

func (c jsonDBConfiguration) GetAccountsSaver() core.DataSaver[[]entities.Account] {
	//TODO implement me
	panic("implement me")
}

func (c jsonDBConfiguration) GetCategoriesSaver() core.DataSaver[[]entities.Category] {
	//TODO implement me
	panic("implement me")
}

func (c jsonDBConfiguration) GetSubcategoriesSaver() core.DataSaver[[]entities.Subcategory] {
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

func (c jsonDBConfiguration) GetSubcategories(fileName string) ([]entities.Subcategory, error) {
	return core.LoadJson[[]entities.Subcategory](fileName + ".json")
}

func (c jsonDBConfiguration) GetMainDataSource() core.DatedSource[entities.FinanceRecord] {
	return jsonDatedSource{}
}
