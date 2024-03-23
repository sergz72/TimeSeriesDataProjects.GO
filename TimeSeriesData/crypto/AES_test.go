package crypto

import (
	"reflect"
	"testing"
)

func TestAESGcm_EncryptWithNonce(t *testing.T) {
	key := make([]byte, 32)
	for i := 0; i < 32; i++ {
		key[i] = byte(i)
	}
	nonce := make([]byte, 12)
	for i := 0; i < 12; i++ {
		nonce[i] = byte(i)
	}
	aes, err := NewAesGcm(key)
	if err != nil {
		t.Fatal(err)
	}
	encrypted := aes.EncryptWithNonce(key, nonce)
	shouldBe := []byte{71, 3, 212, 24, 193, 224, 196, 28, 133, 72, 157, 128, 189, 228, 118, 98, 147, 199, 149, 39, 228, 110,
		73, 107, 32, 126, 255, 158, 1, 116, 30, 173, 94, 221, 220, 80, 116, 4, 78, 34, 130, 180, 50, 179, 242, 216, 246, 115}
	if !reflect.DeepEqual(encrypted, shouldBe) {
		t.Fatal("different encrypted data")
	}
	decrypted, err := aes.DecryptWithNonce(encrypted, nonce)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decrypted, key) {
		t.Fatal("different decrypted data")
	}
}
