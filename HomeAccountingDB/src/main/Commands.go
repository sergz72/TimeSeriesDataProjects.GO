package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type dictsCommand struct{}

func newDictsCommand(buffer *bytes.Buffer) (command, error) {
	if buffer.Len() != 0 {
		return nil, errors.New("invalid dicts command")
	}
	return &dictsCommand{}, nil
}

func (d *dictsCommand) Execute(db *dB) ([]byte, error) {
	return db.getDicts()
}

func (d *dictsCommand) ReadOnlyLockRequired() bool {
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
	return &opsCommand{int(date)}, err
}

func (d *opsCommand) Execute(db *dB) ([]byte, error) {
	return db.getOpsAndChanges(d.date)
}

func (d *opsCommand) ReadOnlyLockRequired() bool {
	return false
}
