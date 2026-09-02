package utils

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"

	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) domain.PasswordHasher {
	return &bcryptHasher{cost: cost}
}

func (b *bcryptHasher) Hash(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *bcryptHasher) Compare(hashedPassword, password []byte) bool {
	if len(hashedPassword) == 0 || len(password) == 0 {
		return false
	}
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err != nil
}
