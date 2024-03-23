package main

import (
	"TimeSeriesData/crypto"
	"bytes"
	"errors"
	"sync"
)

type tcpServerData struct {
	s      settings
	db     *dB
	lock   sync.RWMutex
	aesKey []byte
}

type command interface {
	Execute(db *dB) ([]byte, error)
	ReadOnlyLockRequired() bool
}

func (d *tcpServerData) handle(request []byte) ([]byte, error) {
	if len(request) < 33 {
		return nil, errors.New("too short request")
	}
	cmd, err := decodeRequest(request[32:])
	if err != nil {
		return nil, err
	}
	if d.db == nil {
		d.lock.Lock()
		if d.db == nil {
			aes, err := crypto.NewAesGcm(request[:32])
			if err != nil {
				d.lock.Unlock()
				return nil, err
			}
			c := newBinaryDBConfiguration(aes)
			d.db, err = initDatabase(d.s, c)
			if err != nil {
				d.lock.Unlock()
				return nil, err
			}
		}
		d.lock.Unlock()
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
