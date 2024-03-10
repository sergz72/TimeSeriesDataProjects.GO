package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"os"
)

type AESGcm struct {
	aesgcm cipher.AEAD
}

func (a AESGcm) Encrypt(data []byte) []byte {
	nonce := make([]byte, 12)
	rand.Read(nonce)
	return append(nonce, a.aesgcm.Seal(nil, nonce, data, nil)...)
}

func (a AESGcm) Decrypt(data []byte) ([]byte, error) {
	if len(data) <= 12 {
		return nil, errors.New("wrong data size")
	}
	nonce := data[:12]
	return a.aesgcm.Open(nil, nonce, data[12:], nil)
}

func NewAesGcm(key []byte) (AESGcm, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return AESGcm{}, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return AESGcm{}, err
	}
	return AESGcm{aesgcm: aesgcm}, nil
}

func LoadAesGcmKey(fileName string) ([]byte, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if len(data) != 12 {
		return nil, errors.New("wrong file size")
	}
	return data, nil
}