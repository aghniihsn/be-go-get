package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher interface defines password hashing operations
type Hasher interface {
	Hash(password string) (string, error)
	Verify(password, hashedPassword string) error
}

// BcryptHasher implements Hasher interface using bcrypt
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new bcrypt hasher with specified cost
func NewBcryptHasher(cost int) Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

// NewDefaultHasher creates a bcrypt hasher with default cost
func NewDefaultHasher() Hasher {
	return NewBcryptHasher(bcrypt.DefaultCost)
}

// Hash hashes a plain text password using bcrypt
func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify verifies a password against its hash using bcrypt
func (h *BcryptHasher) Verify(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// HashPassword is a convenience function for quick password hashing
func HashPassword(password string) (string, error) {
	hasher := NewDefaultHasher()
	return hasher.Hash(password)
}

// VerifyPassword is a convenience function for quick password verification
func VerifyPassword(password, hashedPassword string) error {
	hasher := NewDefaultHasher()
	return hasher.Verify(password, hashedPassword)
}
