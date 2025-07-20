package password

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("HashPassword returned empty string")
	}

	if hashedPassword == password {
		t.Fatal("HashPassword returned the same string as input")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Test correct password
	err = VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Fatalf("VerifyPassword failed for correct password: %v", err)
	}

	// Test wrong password
	err = VerifyPassword(wrongPassword, hashedPassword)
	if err == nil {
		t.Fatal("VerifyPassword should fail for wrong password")
	}
}

func TestBcryptHasher(t *testing.T) {
	hasher := NewDefaultHasher()
	password := "testpassword123"

	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("BcryptHasher.Hash failed: %v", err)
	}

	err = hasher.Verify(password, hashedPassword)
	if err != nil {
		t.Fatalf("BcryptHasher.Verify failed for correct password: %v", err)
	}

	err = hasher.Verify("wrongpassword", hashedPassword)
	if err == nil {
		t.Fatal("BcryptHasher.Verify should fail for wrong password")
	}
}

func TestNewBcryptHasherWithCost(t *testing.T) {
	customCost := 12
	hasher := NewBcryptHasher(customCost)
	password := "testpassword123"

	hashedPassword, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("BcryptHasher.Hash with custom cost failed: %v", err)
	}

	err = hasher.Verify(password, hashedPassword)
	if err != nil {
		t.Fatalf("BcryptHasher.Verify with custom cost failed: %v", err)
	}
}

func TestNewBcryptHasherInvalidCost(t *testing.T) {
	// Test with cost too low
	hasher := NewBcryptHasher(1)
	bcryptHasher := hasher.(*BcryptHasher)
	if bcryptHasher.cost != bcrypt.DefaultCost {
		t.Fatalf("Expected default cost %d, got %d", bcrypt.DefaultCost, bcryptHasher.cost)
	}

	// Test with cost too high
	hasher = NewBcryptHasher(100)
	bcryptHasher = hasher.(*BcryptHasher)
	if bcryptHasher.cost != bcrypt.DefaultCost {
		t.Fatalf("Expected default cost %d, got %d", bcrypt.DefaultCost, bcryptHasher.cost)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"
	hasher := NewDefaultHasher()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hasher.Hash(password)
		if err != nil {
			b.Fatalf("Hash failed: %v", err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "testpassword123"
	hasher := NewDefaultHasher()
	hashedPassword, _ := hasher.Hash(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := hasher.Verify(password, hashedPassword)
		if err != nil {
			b.Fatalf("Verify failed: %v", err)
		}
	}
}
