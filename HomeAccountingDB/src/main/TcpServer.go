package main

import (
	"bytes"
	"errors"
	"sync"
)

type tcpServerData struct {
	db     *dB
	lock   *sync.RWMutex
	aesKey []byte
}

type command interface {
	Execute(db *dB) ([]byte, error)
	ReadOnlyLockRequired() bool
}

func (d *tcpServerData) handle(request []byte) ([]byte, error) {
	cmd, err := decodeRequest(request)
	if err != nil {
		return nil, err
	}
	if d.db == nil {
		d.lock = &sync.RWMutex{}
		panic("todo")
	}
	if cmd.ReadOnlyLockRequired() {
		d.lock.RLock()
		defer d.lock.RUnlock()
	} else {
		d.lock.Lock()
		defer d.lock.Unlock()
	}
	return cmd.Execute(d.db)
}

func decodeRequest(request []byte) (command, error) {
	l := len(request)
	if l == 0 {
		return nil, errors.New("too short request")
	}
	buffer := bytes.NewBuffer(request[1:])
	switch request[0] {
	case 0: // DICTS request
		return newDictsCommand(buffer)
	case 1: // DICTS request
		return newOpsCommand(buffer)
	default:
		return nil, errors.New("unknown command")
	}
}
