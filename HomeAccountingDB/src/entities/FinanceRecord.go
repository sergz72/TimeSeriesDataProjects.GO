package entities

import (
	"TimeSeriesData/core"
	"io"
)

type FinanceRecord struct {
	operations []FinanceOperation
	totals     map[int]int
}

type OpsAndChanges struct {
	Operations []FinanceOperation
	Changes    map[int]*FinanceChange
}

func NewFinanceRecord(operations []FinanceOperation) *FinanceRecord {
	return &FinanceRecord{
		operations: operations,
		totals:     make(map[int]int),
	}
}

func CreateChanges(totals map[int]int) map[int]*FinanceChange {
	result := make(map[int]*FinanceChange)
	for k, v := range totals {
		result[k] = &FinanceChange{
			StartBalance:     v,
			SummaIncome:      0,
			SummaExpenditure: 0,
		}
	}
	return result
}

func (r *FinanceRecord) BuildChanges() map[int]*FinanceChange {
	return CreateChanges(r.totals)
}

func BuildTotals(changes map[int]*FinanceChange) map[int]int {
	totals := make(map[int]int)
	for k, v := range changes {
		totals[k] = v.GetEndSumma()
	}
	return totals
}

func (r *FinanceRecord) SetTotals(changes map[int]*FinanceChange) {
	r.totals = BuildTotals(changes)
}

func (r *FinanceRecord) UpdateChanges(changes map[int]*FinanceChange, accounts core.DictionaryData[Account],
	subcategories core.DictionaryData[Subcategory], from, to int) error {
	for _, op := range r.operations {
		if op.Date >= from && op.Date <= to {
			err := op.UpdateChanges(changes, accounts, subcategories)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *FinanceRecord) BuildOpsAndChanges(date int, accounts core.DictionaryData[Account],
	subcategories core.DictionaryData[Subcategory]) (OpsAndChanges, error) {
	changes := r.BuildChanges()
	err := r.UpdateChanges(changes, accounts, subcategories, 0, date-1)
	if err != nil {
		return OpsAndChanges{}, err
	}
	totals := BuildTotals(changes)
	changes = CreateChanges(totals)
	var ops []FinanceOperation
	for _, op := range r.operations {
		if op.Date == date {
			ops = append(ops, op)
		}
	}
	err = r.UpdateChanges(changes, accounts, subcategories, date, date)
	if err != nil {
		return OpsAndChanges{}, err
	}
	return OpsAndChanges{
		Operations: ops,
		Changes:    changes,
	}, nil
}

func (r *FinanceRecord) Save(writer io.Writer) error {
	//TODO implement me
	panic("implement me")
}

func NewFinanceRecordFromBinary(reader io.Reader) (*FinanceRecord, error) {
	return nil, nil
}
