package main

import (
	"HomeAccountingDB/src/entities"
	"TimeSeriesData/core"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type dictsCommand struct{}

func newDictsCommand(buffer *bytes.Buffer) (command, error) {
	if buffer.Len() != 0 {
		return nil, errors.New("invalid dicts command")
	}
	fmt.Println("dicts command")
	return &dictsCommand{}, nil
}

func (c *dictsCommand) Execute(db *dB) ([]byte, error) {
	return db.getDicts()
}

func (c *dictsCommand) ReadOnlyLockRequired() bool {
	return true
}

type opsCommand struct {
	date int
}

func newOpsCommand(buffer *bytes.Buffer) (command, error) {
	if buffer.Len() != 4 {
		return nil, errors.New("invalid ops command")
	}
	var date uint32
	err := binary.Read(buffer, binary.LittleEndian, &date)
	fmt.Printf("ops command, date=%v\n", date)
	return &opsCommand{int(date)}, err
}

func (c *opsCommand) Execute(db *dB) ([]byte, error) {
	return db.getOpsAndChanges(c.date)
}

func (c *opsCommand) ReadOnlyLockRequired() bool {
	return false
}

type addOperationCommand struct {
	date        int
	subcategory int
	account     int
	summa       string
	amount      string
	properties  []entities.FinOpProperty
}

func newAddOperationCommand(buffer *bytes.Buffer, name string) (command, error) {
	if buffer.Len() < 18 {
		return nil, fmt.Errorf("invalid %v command", name)
	}
	var date uint32
	err := binary.Read(buffer, binary.LittleEndian, &date)
	if err != nil {
		return nil, err
	}
	var subcategory uint32
	err = binary.Read(buffer, binary.LittleEndian, &subcategory)
	if err != nil {
		return nil, err
	}
	var account uint32
	err = binary.Read(buffer, binary.LittleEndian, &account)
	if err != nil {
		return nil, err
	}
	summa, err := core.ReadStringFromBinary(buffer)
	if err != nil {
		return nil, err
	}
	amount, err := core.ReadStringFromBinary(buffer)
	if err != nil {
		return nil, err
	}
	var l uint16
	err = binary.Read(buffer, binary.LittleEndian, &l)
	if err != nil {
		return nil, err
	}
	var properties []entities.FinOpProperty
	for l > 0 {
		prop, err := entities.NewFinOpPropertyFromBinary(buffer)
		if err != nil {
			return nil, err
		}
		properties = append(properties, prop)
		l--
	}
	if buffer.Len() > 0 {
		return nil, fmt.Errorf("incorrect %v command length", name)
	}
	c := addOperationCommand{
		date:        int(date),
		subcategory: int(subcategory),
		account:     int(account),
		summa:       summa,
		amount:      amount,
		properties:  nil,
	}
	fmt.Printf("%v command %v\n", name, c)
	return &c, nil
}

func (c *addOperationCommand) Execute(db *dB) ([]byte, error) {
	return db.addOperation(c)
}

func (c *addOperationCommand) ReadOnlyLockRequired() bool {
	return false
}

type modifyOperationCommand addOperationCommand

func newModifyOperationCommand(buffer *bytes.Buffer) (command, error) {
	c, err := newAddOperationCommand(buffer, "modifyOperation")
	ac := c.(*addOperationCommand)
	mc := modifyOperationCommand(*ac)
	return &mc, err
}

func (c *modifyOperationCommand) Execute(db *dB) ([]byte, error) {
	return db.modifyOperation(c)
}

func (c *modifyOperationCommand) ReadOnlyLockRequired() bool {
	return false
}

type deleteOperationCommand struct {
	date        int
	subcategory int
	account     int
}

func newDeleteOperationCommand(buffer *bytes.Buffer) (command, error) {
	if buffer.Len() != 12 {
		return nil, errors.New("invalid deleteOperation command")
	}
	var date uint32
	err := binary.Read(buffer, binary.LittleEndian, &date)
	if err != nil {
		return nil, err
	}
	var subcategory uint32
	err = binary.Read(buffer, binary.LittleEndian, &subcategory)
	if err != nil {
		return nil, err
	}
	var account uint32
	err = binary.Read(buffer, binary.LittleEndian, &account)
	fmt.Printf("deleteOperation command, date=%v subcategory=%v account=%v\n", date, subcategory, account)
	return &deleteOperationCommand{int(date), int(subcategory), int(account)}, err
}

func (c *deleteOperationCommand) Execute(db *dB) ([]byte, error) {
	return db.deleteOperation(c)
}

func (c *deleteOperationCommand) ReadOnlyLockRequired() bool {
	return false
}
