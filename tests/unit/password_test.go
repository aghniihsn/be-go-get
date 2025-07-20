package main_test

import (
	"go-get-backend/pkg/password"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPasswordHashing tests password hashing functionality
func TestPasswordHashing(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Simple password",
			password: "password123",
		},
		{
			name:     "Complex password",
			password: "MyC0mpl3x!P@ssw0rd",
		},
		{
			name:     "Short password",
			password: "abc",
		},
		{
			name:     "Long password",
			password: "ThisIsAVeryLongPasswordThatShouldStillWorkCorrectly123456789",
		},
		{
			name:     "Password with special characters",
			password: "P@ssw0rd!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test hashing
			hashedPassword, err := password.HashPassword(tt.password)
			assert.NoError(t, err, "Password hashing should not return error")
			assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
			assert.NotEqual(t, tt.password, hashedPassword, "Hashed password should differ from original")

			// Test verification with correct password
			err = password.VerifyPassword(tt.password, hashedPassword)
			assert.NoError(t, err, "Password verification should succeed with correct password")

			// Test verification with incorrect password
			err = password.VerifyPassword("wrongpassword", hashedPassword)
			assert.Error(t, err, "Password verification should fail with incorrect password")
		})
	}
}

// TestPasswordConsistency tests that same password produces different hashes
func TestPasswordConsistency(t *testing.T) {
	testPassword := "testpassword123"

	// Generate multiple hashes of the same password
	hashes := make([]string, 5)
	for i := 0; i < 5; i++ {
		hash, err := password.HashPassword(testPassword)
		assert.NoError(t, err)
		hashes[i] = hash
	}

	// Verify all hashes are different (due to salt)
	for i := 0; i < len(hashes); i++ {
		for j := i + 1; j < len(hashes); j++ {
			assert.NotEqual(t, hashes[i], hashes[j], "Each hash should be unique due to salt")
		}
	}

	// Verify all hashes verify correctly
	for _, hash := range hashes {
		err := password.VerifyPassword(testPassword, hash)
		assert.NoError(t, err, "All hashes should verify correctly")
	}
}

// TestPasswordEdgeCases tests edge cases for password functionality
func TestPasswordEdgeCases(t *testing.T) {
	t.Run("Empty password", func(t *testing.T) {
		hashedPassword, err := password.HashPassword("")
		assert.NoError(t, err, "Should handle empty password")

		err = password.VerifyPassword("", hashedPassword)
		assert.NoError(t, err, "Should verify empty password correctly")
	})

	t.Run("Unicode password", func(t *testing.T) {
		unicodePassword := "пароль🔐密码"
		hashedPassword, err := password.HashPassword(unicodePassword)
		assert.NoError(t, err, "Should handle unicode password")

		err = password.VerifyPassword(unicodePassword, hashedPassword)
		assert.NoError(t, err, "Should verify unicode password correctly")
	})

	t.Run("Verify with empty hash", func(t *testing.T) {
		err := password.VerifyPassword("password123", "")
		assert.Error(t, err, "Should fail with empty hash")
	})

	t.Run("Verify with invalid hash format", func(t *testing.T) {
		err := password.VerifyPassword("password123", "invalid_hash_format")
		assert.Error(t, err, "Should fail with invalid hash format")
	})
}

// BenchmarkPasswordHashing benchmarks password hashing performance
func BenchmarkPasswordHashing(b *testing.B) {
	testPassword := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = password.HashPassword(testPassword)
	}
}

// BenchmarkPasswordVerification benchmarks password verification performance
func BenchmarkPasswordVerification(b *testing.B) {
	testPassword := "benchmarkpassword123"
	hashedPassword, _ := password.HashPassword(testPassword)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = password.VerifyPassword(testPassword, hashedPassword)
	}
}
