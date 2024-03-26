package entities

import (
	"TimeSeriesData/core"
	"encoding/binary"
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

func (c OpsAndChanges) Save(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint16(len(c.Operations)))
	if err != nil {
		return err
	}
	for _, op := range c.Operations {
		err = op.SaveToBinary(writer)
		if err != nil {
			return err
		}
	}
	err = binary.Write(writer, binary.LittleEndian, uint16(len(c.Changes)))
	if err != nil {
		return err
	}
	for accountId, change := range c.Changes {
		err = binary.Write(writer, binary.LittleEndian, uint16(accountId))
		if err != nil {
			return err
		}
		err = change.SaveToBinary(writer)
		if err != nil {
			return err
		}
	}
	return nil
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
	subcategories core.DictionaryData[Subcategory], filterChanges bool) (OpsAndChanges, error) {
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
	if filterChanges {
		changes, err = doFilterChanges(changes, accounts, date)
		if err != nil {
			return OpsAndChanges{}, err
		}
	}
	return OpsAndChanges{
		Operations: ops,
		Changes:    changes,
	}, nil
}

func doFilterChanges(changes map[int]*FinanceChange, accounts core.DictionaryData[Account],
	date int) (map[int]*FinanceChange, error) {
	result := make(map[int]*FinanceChange)
	for accountId, change := range changes {
		acc, err := accounts.Get(accountId)
		if err != nil {
			return nil, err
		}
		if acc.ActiveTo > 0 && int(acc.ActiveTo) <= date {
			continue
		}
		result[accountId] = change
	}
	return result, nil
}

func (r *FinanceRecord) Save(writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint16(len(r.operations)))
	if err != nil {
		return err
	}
	for _, op := range r.operations {
		err = op.SaveToBinary(writer)
		if err != nil {
			return err
		}
	}
	err = binary.Write(writer, binary.LittleEndian, uint16(len(r.totals)))
	if err != nil {
		return err
	}
	for k, v := range r.totals {
		err = binary.Write(writer, binary.LittleEndian, uint16(k))
		if err != nil {
			return err
		}
		err = binary.Write(writer, binary.LittleEndian, int64(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *FinanceRecord) BuildHints() map[FinOpPropertyCode]map[string]bool {
	result := make(map[FinOpPropertyCode]map[string]bool)
	for _, op := range r.operations {
		for _, prop := range op.FinOpProperties {
			if prop.StringValue != nil && len(*prop.StringValue) > 0 {
				h, ok := result[prop.PropertyCode]
				if !ok {
					h = make(map[string]bool)
					result[prop.PropertyCode] = h
				}
				h[*prop.StringValue] = true
			}
		}
	}
	return result
}

func (r *FinanceRecord) Copy(from, to int) *FinanceRecord {
	return &FinanceRecord{
		operations: r.GetOperations(from, to),
		totals:     r.totals,
	}
}

func (r *FinanceRecord) GetOperations(from, to int) []FinanceOperation {
	var result []FinanceOperation
	for _, op := range r.operations {
		if op.Date >= from && op.Date <= to {
			result = append(result, op)
		}
	}
	return result
}

func (r *FinanceRecord) AddOperations(operations []FinanceOperation) {
	r.operations = append(r.operations, operations...)
}

func NewFinanceRecordFromBinary(reader io.Reader) (*FinanceRecord, error) {
	var r FinanceRecord
	var l uint16
	err := binary.Read(reader, binary.LittleEndian, &l)
	if err != nil {
		return nil, err
	}
	for l > 0 {
		var op FinanceOperation
		op, err = NewFinanceOperationFromBinary(reader)
		if err != nil {
			return nil, err
		}
		r.operations = append(r.operations, op)
		l--
	}
	err = binary.Read(reader, binary.LittleEndian, &l)
	if err != nil {
		return nil, err
	}
	r.totals = make(map[int]int)
	for l > 0 {
		var k uint16
		err = binary.Read(reader, binary.LittleEndian, &k)
		if err != nil {
			return nil, err
		}
		var v int64
		err = binary.Read(reader, binary.LittleEndian, &v)
		if err != nil {
			return nil, err
		}
		r.totals[int(k)] = int(v)
		l--
	}
	return &r, nil
}
