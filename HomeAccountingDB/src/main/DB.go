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
	Key                    string
}

type dBConfiguration interface {
	GetAccounts(fileName string) ([]entities.Account, error)
	GetCategories(fileName string) ([]entities.Category, error)
	GetSubcategories(fileName, mapFileName string) ([]entities.Subcategory, error)
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

func getSubcategoriesMapFileName(dataFolderPath string) string {
	return dataFolderPath + "/subcategory_to_property_code_map"
}

func getMainDataFolderPath(dataFolderPath string) string {
	return dataFolderPath + "/dates"
}

func indexCalculator(date int, minYear int, minMonth int) int {
	date /= 100
	year := date / 100
	month := date % 100
	return month - minMonth + (year-minYear)*12
}

func loadDicts(s settings, configuration dBConfiguration) (core.DictionaryData[entities.Account],
	core.DictionaryData[entities.Category], core.DictionaryData[entities.Subcategory], error) {
	path := getAccountsFileName(s.DataFolderPath)
	alist, err := configuration.GetAccounts(path)
	if err != nil {
		return core.DictionaryData[entities.Account]{}, core.DictionaryData[entities.Category]{},
			core.DictionaryData[entities.Subcategory]{}, err
	}
	accounts := core.NewDictionaryData[entities.Account](path, "account", alist)
	path = getCategoriesFileName(s.DataFolderPath)
	clist, err := configuration.GetCategories(path)
	if err != nil {
		return core.DictionaryData[entities.Account]{}, core.DictionaryData[entities.Category]{},
			core.DictionaryData[entities.Subcategory]{}, err
	}
	categories := core.NewDictionaryData[entities.Category](path, "category", clist)
	path = getSubcategoriesFileName(s.DataFolderPath)
	slist, err := configuration.GetSubcategories(path, getSubcategoriesMapFileName(s.DataFolderPath))
	if err != nil {
		return core.DictionaryData[entities.Account]{}, core.DictionaryData[entities.Category]{},
			core.DictionaryData[entities.Subcategory]{}, err
	}
	subcategories := core.NewDictionaryData[entities.Subcategory](path, "subcategory", slist)
	return accounts, categories, subcategories, nil
}

func loadDB(s settings, configuration dBConfiguration) (*dB, error) {
	accounts, categories, subcategories, err := loadDicts(s, configuration)
	if err != nil {
		return nil, err
	}
	data, err := core.LoadTimeSeriesData[entities.FinanceRecord](getMainDataFolderPath(s.DataFolderPath),
		configuration.GetMainDataSource(), s.TimeSeriesDataCapacity, func(date int) int {
			return indexCalculator(date, s.MinYear, s.MinMonth)
		}, func(date int) int {
			return date / 100
		}, 1000000)
	if err != nil {
		return nil, err
	}
	return &dB{dataFolderPath: s.DataFolderPath, configuration: configuration, accounts: accounts,
		categories: categories, subcategories: subcategories, data: data}, nil
}

func initDB(s settings, configuration dBConfiguration) (*dB, error) {
	accounts, categories, subcategories, err := loadDicts(s, configuration)
	if err != nil {
		return nil, err
	}
	data, err := core.InitTimeSeriesData[entities.FinanceRecord](getMainDataFolderPath(s.DataFolderPath),
		configuration.GetMainDataSource(), s.TimeSeriesDataCapacity, func(date int) int {
			return indexCalculator(date, s.MinYear, s.MinMonth)
		}, func(date int) int {
			return date / 100
		}, 1000000)
	if err != nil {
		return nil, err
	}
	return &dB{dataFolderPath: s.DataFolderPath, configuration: configuration, accounts: accounts,
		categories: categories, subcategories: subcategories, data: data}, nil
}

func (d *dB) buildTotals(from int) error {
	i, err := d.data.Iterator(from, 99999999)
	if err != nil {
		return err
	}
	var changes map[int]*entities.FinanceChange
	for i.HasNext() {
		idx, v, err := i.Next()
		if err != nil {
			return err
		}
		if changes == nil {
			changes = v.BuildChanges()
		} else {
			v.SetTotals(changes)
			d.data.MarkAsModified(idx)
		}
		err = v.UpdateChanges(changes, d.accounts, d.subcategories, 0, 99999999)
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
		fmt.Println(d.data.GetDate(key))
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
	err := d.accounts.SaveTo(configuration.GetAccountsSaver(), getAccountsFileName(dataFolderPath), entities.SaveAccountByIndex)
	if err != nil {
		return err
	}
	err = d.categories.SaveTo(configuration.GetCategoriesSaver(), getCategoriesFileName(dataFolderPath), entities.SaveCategoryByIndex)
	if err != nil {
		return err
	}
	err = d.subcategories.SaveTo(configuration.GetSubcategoriesSaver(), getSubcategoriesFileName(dataFolderPath), entities.SaveSubcategoryByIndex)
	if err != nil {
		return err
	}
	return d.data.SaveAll(configuration.GetMainDataSource(), getMainDataFolderPath(dataFolderPath))
}

func (d *dB) getDicts() []byte {
	return nil
}
