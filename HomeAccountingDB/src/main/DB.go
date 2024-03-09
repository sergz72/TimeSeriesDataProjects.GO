package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"fmt"
)

type DBConfiguration interface {
	GetAccounts(fileName string) ([]entities.Account, error)
	GetCategories(fileName string) ([]entities.Category, error)
	GetSubcategories(fileName string) ([]entities.Subcategory, error)
	GetMainDataSource() core.DatedSource[entities.FinanceRecord]
}

type DB struct {
	dataFolderPath string
	accounts       core.DictionaryData[entities.Account]
	categories     core.DictionaryData[entities.Category]
	subcategories  core.DictionaryData[entities.Subcategory]
	data           core.TimeSeriesData[entities.FinanceRecord]
}

func LoadDB(minYear int, minMonth int, dataFolderPath string, configuration DBConfiguration, timeSeriesDataCapacity int) (*DB, error) {
	path := dataFolderPath + "/accounts"
	alist, err := configuration.GetAccounts(path)
	if err != nil {
		return nil, err
	}
	accounts := core.NewDictionaryData[entities.Account](path, "account", alist)
	path = dataFolderPath + "/categories"
	clist, err := configuration.GetCategories(path)
	if err != nil {
		return nil, err
	}
	categories := core.NewDictionaryData[entities.Category](path, "category", clist)
	path = dataFolderPath + "/subcategories"
	slist, err := configuration.GetSubcategories(path)
	if err != nil {
		return nil, err
	}
	subcategories := core.NewDictionaryData[entities.Subcategory](path, "subcategory", slist)
	data, err := core.LoadTimeSeriesData[entities.FinanceRecord](dataFolderPath+"/dates",
		configuration.GetMainDataSource(), timeSeriesDataCapacity, func(date int) int {
			date /= 100
			year := date / 100
			month := date % 100
			return month - minMonth + (year-minYear)*12
		}, func(idx int) int {
			month := idx + minMonth
			year := minYear + month/12
			return year*100 + (month % 12)
		}, 1000000)
	if err != nil {
		return nil, err
	}
	return &DB{accounts: accounts, categories: categories, subcategories: subcategories, data: data}, nil
}

func (d *DB) BuildTotals(from int) error {
	m, err := d.data.GetRange(from, 99999999)
	if err != nil {
		return err
	}
	var changes map[int]*entities.FinanceChange
	for _, v := range m {
		if changes == nil {
			changes = v.Data.BuildChanges()
		} else {
			v.Data.SetTotals(changes)
		}
		err = v.Data.UpdateChanges(changes, d.accounts, d.subcategories, 0, 99999999)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) PrintChanges(date int) {
	key, v, err := d.data.Get(date)
	if err != nil {
		panic(err)
	}
	if v != nil {
		var result entities.OpsAndChanges
		result, err = v.BuildOpsAndChanges(date, d.accounts, d.subcategories)
		if err != nil {
			panic(err)
		}
		fmt.Println(d.data.DateCalculator(key))
		for acc, ch := range result.Changes {
			var account *entities.Account
			account, err = d.accounts.Get(acc)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v %v %v %v %v\n", account.GetName(), ch.StartBalance,
				ch.SummaIncome, ch.SummaExpenditure, ch.GetEndSumma())
		}
	} else {
		fmt.Println("no data")
	}
}
