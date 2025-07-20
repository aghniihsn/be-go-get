package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	fmt.Println("🔐 JWT Secret Generator")
	fmt.Println("======================")

	// Generate 32 bytes (256 bits) random data
	secretBytes := make([]byte, 32)
	_, err := rand.Read(secretBytes)
	if err != nil {
		log.Fatal("Failed to generate random bytes:", err)
	}

	// Option 1: Base64 encoded (recommended for JWT)
	base64Secret := base64.URLEncoding.EncodeToString(secretBytes)
	fmt.Printf("Base64 Secret (Recommended):\n%s\n\n", base64Secret)

	// Option 2: Hex encoded
	hexSecret := hex.EncodeToString(secretBytes)
	fmt.Printf("Hex Secret:\n%s\n\n", hexSecret)

	// Option 3: Generate multiple secrets
	fmt.Println("Multiple Secrets (choose one):")
	fmt.Println("------------------------------")
	for i := 1; i <= 3; i++ {
		secretBytes := make([]byte, 32)
		rand.Read(secretBytes)
		secret := base64.URLEncoding.EncodeToString(secretBytes)
		fmt.Printf("%d. %s\n", i, secret)
	}

	fmt.Println("\n💡 Copy one of the Base64 secrets to your .env file")
	fmt.Println("💡 Each secret is 256-bit (32 bytes) for maximum security")
}
