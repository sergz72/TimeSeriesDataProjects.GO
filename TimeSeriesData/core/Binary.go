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

type BinarySaver[T BinaryData] struct {
	processor CryptoProcessor
}

func (b BinarySaver[T]) Save(data T, fileName string) error {
	return SaveBinary(fileName, b.processor, data)
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