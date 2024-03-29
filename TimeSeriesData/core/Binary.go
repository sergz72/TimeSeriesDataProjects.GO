package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"reflect"
)

type CryptoProcessor interface {
	Encrypt(data []byte) []byte
	Decrypt(data []byte) ([]byte, error)
}

type BinaryData interface {
	Save(writer io.Writer) error
}

type BinarySaver struct {
	processor CryptoProcessor
	data *bytes.Buffer
}

func NewBinarySaver(processor CryptoProcessor) *BinarySaver {
	return &BinarySaver{processor: processor, data: new(bytes.Buffer)}
}

func (b *BinarySaver) Save(data any, saveIndex func(int, any, io.Writer) error) error {
	bdata, ok := data.(BinaryData)
	if ok {
		return bdata.Save(b.data)
	}
	t := reflect.ValueOf(data)
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		l := t.Len()
		err := binary.Write(b.data, binary.LittleEndian, uint16(l))
		if err != nil {
			return err
		}
		for i := 0; i < l; i++ {
			err = saveIndex(i, data, b.data)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("unsupported data type")
}

func (b *BinarySaver) GetBytes() []byte {
	dataBytes := b.data.Bytes()
	if b.processor != nil {
		dataBytes = b.processor.Encrypt(dataBytes)
	}
	return dataBytes
}

func (b *BinarySaver) GetFileExtension() string {
	return ".bin"
}

func LoadBinaryArray[T any](reader io.Reader, creator func(reader io.Reader) (T, error)) ([]T, error) {
	var l uint16
	err := binary.Read(reader, binary.LittleEndian, &l)
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

func LoadBinary[T any](fileName string, processor CryptoProcessor, creator func(reader io.Reader) (T, error)) (T, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		var object T
		return object, err
	}
	return LoadBinaryData(data, processor, creator)
}

func LoadBinaryData[T any](data []byte, processor CryptoProcessor, creator func(reader io.Reader) (T, error)) (T, error) {
	if processor != nil {
		var err error
		data, err = processor.Decrypt(data)
		if err != nil {
			var object T
			return object, err
		}
	}
	buffer := bytes.NewBuffer(data)
	value, err := creator(buffer)
	if err != nil {
		return value, err
	}
	if buffer.Len() != 0 {
		return value, errors.New("non zero buffer length")
	}
	return value, err
}

func LoadBinaryP[T any](fileName string, processor CryptoProcessor, creator func(reader io.Reader) (*T, error)) (*T, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return LoadBinaryDataP(data, processor, creator)
}

func LoadBinaryDataP[T any](data []byte, processor CryptoProcessor, creator func(reader io.Reader) (*T, error)) (*T, error) {
	if processor != nil {
		var err error
		data, err = processor.Decrypt(data)
		if err != nil {
			return nil, err
		}
	}
	buffer := bytes.NewBuffer(data)
	value, err := creator(buffer)
	if err != nil {
		return nil, err
	}
	if buffer.Len() != 0 {
		return value, errors.New("non zero buffer length")
	}
	return value, err
}

func buildBytes(processor CryptoProcessor, buffer *bytes.Buffer) []byte {
	data := buffer.Bytes()
	if processor != nil {
		return processor.Encrypt(data)
	}
	return data
}

func buildBinaryDataBytes(processor CryptoProcessor, object BinaryData) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := object.Save(buffer)
	if err != nil {
		return nil, err
	}
	return buildBytes(processor, buffer), nil
}

func SaveBinary(fileName string, processor CryptoProcessor, object BinaryData) error {
	data, err := buildBinaryDataBytes(processor, object)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func ReadStringFromBinary(reader io.Reader) (string, error) {
	var l uint16
	err := binary.Read(reader, binary.LittleEndian, &l)
	if err != nil || l == 0 {
		return "", err
	}
	b := make([]byte, int(l))
	err = binary.Read(reader, binary.LittleEndian, &b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func WriteStringToBinary(writer io.Writer, value string) error {
	var l uint16 = uint16(len(value))
	err := binary.Write(writer, binary.LittleEndian, l)
	if err != nil {
		return err
	}
	if l > 0 {
		return binary.Write(writer, binary.LittleEndian, []byte(value))
	}
	return nil
}
