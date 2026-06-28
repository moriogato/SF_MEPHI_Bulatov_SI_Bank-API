package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// Глобальный ключ для AES (в реальном проекте брать из .env)
var encryptionKey = []byte("32-byte-long-key-for-aes-256-ok!")

func init() {
	// Для AES-256 нужен ключ 32 байта
	if len(encryptionKey) != 32 {
		panic("encryption key must be 32 bytes for AES-256")
	}
}

// EncryptAES шифрует данные с помощью AES-GCM
func EncryptAES(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES расшифровывает данные
func DecryptAES(ciphertextB64 string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashCVV хеширует CVV через bcrypt
func HashCVV(cvv string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckCVV проверяет CVV
func CheckCVV(cvv string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(cvv))
}

// ComputeHMAC вычисляет HMAC-SHA256
func ComputeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC проверяет HMAC
func VerifyHMAC(data, signature string, secret []byte) bool {
	expected := ComputeHMAC(data, secret)
	return expected == signature
}
