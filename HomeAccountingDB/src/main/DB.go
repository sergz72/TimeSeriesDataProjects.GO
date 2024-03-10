package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"fmt"
)

type settings struct {
	MinYear                int
	MinMonth               int
	TimeSeriesDataCapacity int
	DataFolderPath         string
	ServerPort             int
}

type dBConfiguration interface {
	GetAccounts(fileName string) ([]entities.Account, error)
	GetCategories(fileName string) ([]entities.Category, error)
	GetSubcategories(fileName string) ([]entities.Subcategory, error)
	GetMainDataSource() core.DatedSource[entities.FinanceRecord]
	GetAccountsSaver() core.DataSaver[[]entities.Account]
	GetCategoriesSaver() core.DataSaver[[]entities.Category]
	GetSubcategoriesSaver() core.DataSaver[[]entities.Subcategory]
}

type dB struct {
	dataFolderPath string
	configuration  dBConfiguration
	accounts       core.DictionaryData[entities.Account]
	categories     core.DictionaryData[entities.Category]
	subcategories  core.DictionaryData[entities.Subcategory]
	data           core.TimeSeriesData[entities.FinanceRecord]
}

func getAccountsFileName(dataFolderPath string) string {
	return dataFolderPath + "/accounts"
}

func getCategoriesFileName(dataFolderPath string) string {
	return dataFolderPath + "/categories"
}

func getSubcategoriesFileName(dataFolderPath string) string {
	return dataFolderPath + "/subcategories"
}

func getMainDataFolderPath(dataFolderPath string) string {
	return dataFolderPath + "/dates"
}

func loadDB(s settings, configuration dBConfiguration) (*dB, error) {
	path := getAccountsFileName(s.DataFolderPath)
	alist, err := configuration.GetAccounts(path)
	if err != nil {
		return nil, err
	}
	accounts := core.NewDictionaryData[entities.Account](path, "account", alist)
	path = getCategoriesFileName(s.DataFolderPath)
	clist, err := configuration.GetCategories(path)
	if err != nil {
		return nil, err
	}
	categories := core.NewDictionaryData[entities.Category](path, "category", clist)
	path = getSubcategoriesFileName(s.DataFolderPath)
	slist, err := configuration.GetSubcategories(path)
	if err != nil {
		return nil, err
	}
	subcategories := core.NewDictionaryData[entities.Subcategory](path, "subcategory", slist)
	data, err := core.LoadTimeSeriesData[entities.FinanceRecord](getMainDataFolderPath(s.DataFolderPath),
		configuration.GetMainDataSource(), s.TimeSeriesDataCapacity, func(date int) int {
			date /= 100
			year := date / 100
			month := date % 100
			return month - s.MinMonth + (year-s.MinYear)*12
		}, func(idx int) int {
			month := idx + s.MinMonth
			year := s.MinYear + month/12
			return year*100 + (month % 12)
		}, 1000000)
	if err != nil {
		return nil, err
	}
	return &dB{dataFolderPath: s.DataFolderPath, configuration: configuration, accounts: accounts,
		categories: categories, subcategories: subcategories, data: data}, nil
}

func (d *dB) buildTotals(from int) error {
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

func (d *dB) printChanges(date int) {
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

func (d *dB) save() error {
	return d.saveTo(d.dataFolderPath, d.configuration)
}

func (d *dB) saveTo(dataFolderPath string, configuration dBConfiguration) error {
	err := d.accounts.SaveTo(configuration.GetAccountsSaver(), getAccountsFileName(dataFolderPath))
	if err != nil {
		return err
	}
	err = d.categories.SaveTo(configuration.GetCategoriesSaver(), getCategoriesFileName(dataFolderPath))
	if err != nil {
		return err
	}
	err = d.subcategories.SaveTo(configuration.GetSubcategoriesSaver(), getSubcategoriesFileName(dataFolderPath))
	if err != nil {
		return err
	}
	err = d.data.SaveAll(configuration.GetMainDataSource(), getMainDataFolderPath(dataFolderPath))
	return nil
}
