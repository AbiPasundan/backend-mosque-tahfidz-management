package utils

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(append(salt, hash...)), nil
}

func ComparePassword(encodedHash string, password string) (bool, error) {
	data, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}
	if len(data) < 16 {
		return false, nil
	}
	salt := data[:16]
	storedHash := data[16:]
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return string(newHash) == string(storedHash), nil
}
