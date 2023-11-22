package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashRegisterPassword(password string) (string, error) {
	hashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashByte), nil
}

func CompareLoginPasswordAndHash(reqPassword, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(reqPassword))
	if err != nil {
		return err
	}
	return nil
}
