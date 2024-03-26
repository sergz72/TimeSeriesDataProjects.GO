package main

import (
	"TimeSeriesData/crypto"
	"bytes"
	"errors"
	"reflect"
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

func (d *tcpServerData) handle(request []byte) ([]byte, error, bool) {
	if len(request) < 33 {
		// terminate program
		return nil, errors.New("too short request"), true
	}
	aesKey := request[:32]
	if d.aesKey != nil && !reflect.DeepEqual(d.aesKey, aesKey) {
		// terminate program
		return nil, errors.New("wrong AES key"), true
	}
	cmd, err := decodeRequest(request[32:])
	if err != nil {
		return nil, err, false
	}
	if d.db == nil {
		d.lock.Lock()
		if d.db == nil {
			d.aesKey = make([]byte, 32)
			copy(d.aesKey, aesKey)
			aes, err := crypto.NewAesGcm(d.aesKey)
			if err != nil {
				d.lock.Unlock()
				// terminate program
				return nil, err, true
			}
			c := newBinaryDBConfiguration(aes)
			d.db, err = initDatabase(d.s, c)
			if err != nil {
				d.lock.Unlock()
				// terminate program
				return nil, err, true
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
	var data []byte
	data, err = cmd.Execute(d.db)
	return data, err, false
}

func decodeRequest(request []byte) (command, error) {
	buffer := bytes.NewBuffer(request[1:])
	switch request[0] {
	case 0: // DICTS request
		return newDictsCommand(buffer)
	case 1: // DICTS request
		return newOpsCommand(buffer)
	case 2: // DICTS request
		return newOpsRangeCommand(buffer)
	case 3: // DICTS request
		return newAddOperationCommand(buffer, "addOperation")
	case 4: // DICTS request
		return newModifyOperationCommand(buffer)
	case 5: // DICTS request
		return newDeleteOperationCommand(buffer)
	default:
		return nil, errors.New("unknown command")
	}
}
