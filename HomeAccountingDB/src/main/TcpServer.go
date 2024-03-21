package main

import (
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
	switch request[0] {
	case 0: // DICTS request
		return &dictsCommand{}, nil
	}
	return nil, nil
}

type dictsCommand struct{}

func (d *dictsCommand) Execute(db *dB) ([]byte, error) {
	return db.getDicts(), nil
}

func (d *dictsCommand) ReadOnlyLockRequired() bool {
	return true
}
