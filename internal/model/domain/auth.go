package domain

type PasswordHasher interface {
	Hash(password []byte) (string, error)
	Compare(hashedPassword, password []byte) bool
}

type TokenGenerator interface {
	GenerateToken(userID string) (string error)
}
