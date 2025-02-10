package services

import (
	"crypto/rand"
	"encoding/base64"
	"goPocDemo/model"
	"testing"
)

func TestEncrypt(t *testing.T) {
	secretKey := make([]byte, 16)
	if _, err := rand.Read(secretKey); err != nil {
		t.Fatalf("Failed to generate secret key: %v", err)
	}

	tests := []struct {
		name    string
		val     []byte
		secret  []byte
		wantErr bool
	}{
		{
			name:    "Valid encryption",
			val:     []byte("test data"),
			secret:  secretKey,
			wantErr: false,
		},
		{
			name:    "Invalid secret key length",
			val:     []byte("test data"),
			secret:  []byte("short"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Encrypt(tt.val, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	secretKey := make([]byte, 16)
	if _, err := rand.Read(secretKey); err != nil {
		t.Fatalf("Failed to generate secret key: %v", err)
	}

	validData, err := Encrypt([]byte("test data"), secretKey)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	tests := []struct {
		name    string
		val     []byte
		secret  []byte
		wantErr bool
	}{
		{
			name:    "Valid decryption",
			val:     validData,
			secret:  secretKey,
			wantErr: false,
		},
		{
			name:    "Invalid secret key length",
			val:     validData,
			secret:  []byte("short"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.val, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptData(t *testing.T) {
	req := &model.Request{
		Details: map[string]string{
			"card_number": "1234567890123456",
		},
	}

	err := EncryptData(req)
	if err != nil {
		t.Fatalf("EncryptData() error = %v", err)
	}

	encryptedCardNumber := req.Details["card_number"]
	decoded, err := base64.StdEncoding.DecodeString(encryptedCardNumber)
	if err != nil {
		t.Fatalf("Failed to decode encrypted card number: %v", err)
	}

	secretKey := GetSecretKey()
	decrypted, err := Decrypt(decoded, secretKey)
	if err != nil {
		t.Fatalf("Failed to decrypt card number: %v", err)
	}

	if string(decrypted) != "1234567890123456" {
		t.Errorf("Decrypted card number = %v, want %v", string(decrypted), "1234567890123456")
	}
}
