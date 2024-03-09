package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"errors"
	"strconv"
)

type JsonDatedSource struct{}

func (s JsonDatedSource) GetFileDate(_ string, folderName string) (int, error) {
	return strconv.Atoi(folderName)
}

func (s JsonDatedSource) Load(files []core.FileWithDate) (*entities.FinanceRecord, error) {
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

func (s JsonDatedSource) GetFiles(date int) ([]core.FileWithDate, error) {
	return nil, errors.New("not implemented")
}

func (s JsonDatedSource) Save(date int, data *entities.FinanceRecord) error {
	return errors.New("not implemented")
}

type JsonDBConfiguration struct{}

func (c JsonDBConfiguration) GetAccounts(fileName string) ([]entities.Account, error) {
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

func (c JsonDBConfiguration) GetCategories(fileName string) ([]entities.Category, error) {
	return core.LoadJson[[]entities.Category](fileName + ".json")
}

func (c JsonDBConfiguration) GetSubcategories(fileName string) ([]entities.Subcategory, error) {
	return core.LoadJson[[]entities.Subcategory](fileName + ".json")
}

func (c JsonDBConfiguration) GetMainDataSource() core.DatedSource[entities.FinanceRecord] {
	return JsonDatedSource{}
}
