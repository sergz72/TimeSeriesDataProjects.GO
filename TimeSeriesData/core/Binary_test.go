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
	return binary.Write(writer, binary.LittleEndian, uint32(t.Id))
}

func newTestBinaryData(reader io.Reader) (testBinaryData, error) {
	var id uint32
	err := binary.Read(reader, binary.LittleEndian, &id)
	return testBinaryData{int(id)}, err
}

func TestBinarySaver(t *testing.T) {
	source := testBinaryData{1}
	data, err := BinarySaver[testBinaryData]{}.BuildBytes(source, nil)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadBinaryData[testBinaryData](data, nil, newTestBinaryData)
	if !reflect.DeepEqual(source, loaded) {
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

	data, err = NewBinarySaver[testBinaryData](processor).BuildBytes(source, nil)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err = LoadBinaryData[testBinaryData](data, processor, newTestBinaryData)
	if !reflect.DeepEqual(source, loaded) {
		t.Fatal("different data2")
	}
}

func saveIndex(index int, value []testBinaryData, writer io.Writer) error {
	return value[index].Save(writer)
}

func TestBinarySaverArray(t *testing.T) {
	source := []testBinaryData{{1}, {2}, {3}}
	data, err := BinarySaver[[]testBinaryData]{}.BuildBytes(source, saveIndex)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadBinaryData[[]testBinaryData](data, nil, func(reader io.Reader) ([]testBinaryData, error) {
		return LoadBinaryArray(reader, newTestBinaryData)
	})
	if len(loaded) != 3 || loaded[0].Id != 1 || loaded[1].Id != 2 || loaded[2].Id != 3 {
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

	data, err = NewBinarySaver[[]testBinaryData](processor).BuildBytes(source, saveIndex)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err = LoadBinaryData[[]testBinaryData](data, processor, func(reader io.Reader) ([]testBinaryData, error) {
		return LoadBinaryArray(reader, newTestBinaryData)
	})
	if len(loaded) != 3 || loaded[0].Id != 1 || loaded[1].Id != 2 || loaded[2].Id != 3 {
		t.Fatal("different data2")
	}
}
