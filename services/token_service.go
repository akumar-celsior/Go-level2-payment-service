package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"goPocDemo/model"
	"io"
	"log"
	"sync"
)

type TokenService interface {
	EncryptData(req *model.Request) error
}

var secretKey []byte
var once sync.Once

func Decrypt(val []byte, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	size := aead.NonceSize()
	if len(val) < size {
		return nil, err
	}

	result, err := aead.Open(nil, val[:size], val[size:], nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Encrypt(val []byte, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return aead.Seal(nonce, nonce, val, nil), nil
}

func EncryptData(req *model.Request) error {
	secretKey := GetSecretKey()
	encrypted, err := Encrypt([]byte(req.Details["card_number"]), secretKey)
	if err != nil {
		return errors.New("Unable to encrypt the data")
	}
	req.Details["card_number"] = base64.StdEncoding.EncodeToString(encrypted)
	// fmt.Println(req.Details["card_number"], encrypted, secretKey)
	// decrypted, err := Decrypt(encrypted, secretKey)
	// fmt.Println("first check", encrypted, decrypted, secretKey)
	return nil
}

// GetSecretKey returns the secret key, generating it only once.
func GetSecretKey() []byte {
	// Ensure the key is generated only once
	once.Do(func() {
		var err error
		secretKey, err = generateSecretKey()
		if err != nil {
			log.Fatalf("Error generating secret key: %v", err)
		}
	})
	return secretKey
}

func generateSecretKey() ([]byte, error) {
	key := make([]byte, 16) // 16 bytes = 128 bits
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}
