package core

import (
	"TimeSeriesData/crypto"
	"crypto/rand"
	"encoding/binary"
	"io"
	"reflect"
	"testing"
)

type testBinaryData struct {
	Id int
}

func (t testBinaryData) Save(writer io.Writer) error {
	return binary.Write(writer, binary.BigEndian, uint32(t.Id))
}

func newTestBinaryData(reader io.Reader) (testBinaryData, error) {
	var id uint32
	err := binary.Read(reader, binary.BigEndian, &id)
	return testBinaryData{int(id)}, err
}

func TestBinarySaver(t *testing.T) {
	data, err := BinarySaver[testBinaryData]{}.buildBytes(testBinaryData{1})
	if err != nil {
		t.Fatal(err)
	}
	data2, err := LoadBinaryData[testBinaryData](data, nil, newTestBinaryData)
	if reflect.DeepEqual(data, data2) {
		t.Fatal("different data")
	}
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		t.Fatal(err)
	}
	processor, err := crypto.NewAesGcm(key)
	if err != nil {
		t.Fatal(err)
	}

	data, err = NewBinarySaver[testBinaryData](processor).buildBytes(testBinaryData{1})
	if err != nil {
		t.Fatal(err)
	}
	data2, err = LoadBinaryData[testBinaryData](data, processor, newTestBinaryData)
	if reflect.DeepEqual(data, data2) {
		t.Fatal("different data2")
	}
}
