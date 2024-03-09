package core

import "io"

type CryptoProcessor interface {
	Encode(data []byte)
	Decode(data []byte)
}

type BinaryData interface {
	Save(writer io.Writer)
}

func LoadBinary[T BinaryData](fileName string, processor CryptoProcessor, creator func(reader io.Reader) *T) (*T, error) {
	return nil, nil
}