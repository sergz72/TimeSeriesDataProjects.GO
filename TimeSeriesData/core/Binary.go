package core

import (
	"encoding/binary"
	"errors"
	"bytes"
	"io"
	"os"
)

type CryptoProcessor interface {
	Encrypt(data []byte) []byte
	Decrypt(data []byte) ([]byte, error)
}

type BinaryData interface {
	Save(writer io.Writer) error
}

type BinarySaver[T any] struct {
	processor CryptoProcessor
}

func (b BinarySaver[T]) Save(data T, fileName string) error {
	bdata, ok := any(data).(BinaryData)
	if ok {
		return SaveBinary(fileName, b.processor, bdata)
	}
	adata, ok := any(data).([]BinaryData)
	if ok {
		buffer := new(bytes.Buffer)
		err := binary.Write(buffer, binary.BigEndian, uint16(len(adata)))
		if err != nil {
			return err
		}
		for _, item := range adata {
			err = item.Save(buffer)
			if err != nil {
				return err
			}
		}
		return SaveBinaryBuffer(fileName, b.processor, buffer)
	}
	return errors.New("unsupported data type")
}

func NewBinaryData[T any](reader io.Reader, creator func(reader io.Reader) (T, error)) ([]T, error) {
	var l uint16
	err := binary.Read(reader, binary.BigEndian, &l)
	if err != nil {
		return nil, err
	}
	var result []T
	for l > 0 {
		v, err := creator(reader)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
		l--
	}
	return result, nil
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

func SaveBinaryBuffer(fileName string, processor CryptoProcessor, buffer *bytes.Buffer) error {
	data := buffer.Bytes()
	if processor != nil {
		var err error
		data = processor.Encrypt(data)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(fileName, data, 0644)
}

func SaveBinary(fileName string, processor CryptoProcessor, object BinaryData) error {
	buffer := new(bytes.Buffer)
	err := object.Save(buffer)
	if err != nil {
		return err
	}
	return SaveBinaryBuffer(fileName, processor, buffer)
}

func ReadStringFromBinary(reader io.Reader) (string, error) {
	var l uint16
	err := binary.Read(reader, binary.BigEndian, &l)
	if err != nil || l == 0 {
		return "", err
	}
	b := make([]byte, int(l))
	err = binary.Read(reader, binary.BigEndian, &b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WriteStringToBinary(writer io.Writer, value string) error {
	var l uint16 = uint16(len(value))
	err := binary.Write(writer, binary.BigEndian, l)
	if err != nil {
		return err
	}
	if l > 0 {
		return binary.Write(writer, binary.BigEndian, []byte(value))
	}
	return nil
}