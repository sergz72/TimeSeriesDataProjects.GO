package core

import (
	"encoding/binary"
	"bytes"
	"io"
	"os"
)

type CryptoProcessor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type BinaryData interface {
	Save(writer io.Writer) error
}

func CreateFromBinary[T any](data io.Reader) (T, error){
	var object T
	err := binary.Read(data, binary.BigEndian, &object)
	return object, err
}

func LoadBinary[T any](fileName string, processor CryptoProcessor,
	creator func(reader io.Reader) (T, error)) (T, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		var object T
		return object, err
	}
	if processor != nil {
		data, err = processor.Decrypt(data)
		if err != nil {
			var object T
			return object, err
		}
	}
	return creator(bytes.NewBuffer(data))
}

func LoadBinaryP[T any](fileName string, processor CryptoProcessor,
	creator func(reader io.Reader) (*T, error)) (*T, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if processor != nil {
		data, err = processor.Decrypt(data)
		if err != nil {
			return nil, err
		}
	}
	return creator(bytes.NewBuffer(data))
}

func SaveBinary(fileName string, processor CryptoProcessor, object BinaryData) error {
	buffer := new(bytes.Buffer)
	err := object.Save(buffer)
	if err != nil {
		return err
	}
	data := buffer.Bytes()
	if processor != nil {
		data, err = processor.Encrypt(data)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(fileName, data, 0644)
}