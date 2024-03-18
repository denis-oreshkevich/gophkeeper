package auth

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordMismatch = errors.New("passwords mismatch")

func EncryptPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd, plainPwd string) error {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return ErrPasswordMismatch
	}
	return nil
}
