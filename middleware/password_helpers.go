package middleware

import (
	"encoding/hex"
	"errors"
	"math/rand"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, string, error) {

	s := make([]byte, 16)

	_, err := rand.Read(s)
	if err != nil {
		return "", "", err
	}
	salt := hex.EncodeToString(s)

	h := argon2.IDKey([]byte(password), s, 1, 64*1024, 4, 32)
	hash := hex.EncodeToString(h)

	return salt, hash, nil
}

func VerifyPassword(password string, salt string, hash string) error {

	data, err := hex.DecodeString(salt)
	if err != nil {
		return err
	}

	h := argon2.IDKey([]byte(password), data, 1, 64*1024, 4, 32)
	generatedHash := hex.EncodeToString(h)

	if generatedHash != hash {
		return errors.New("wrong_password")
	}

	return nil
}
