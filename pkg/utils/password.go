package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func ComparePassword(hashedPassword []byte, plainPassword []byte) bool {
	if len(hashedPassword) == 0 || len(plainPassword) == 0 {
		return false
	}
	err := bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)
	return err != nil
}
