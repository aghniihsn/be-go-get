package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("🔐 JWT Secret Auto-Updater")
	fmt.Println("===========================")

	// Generate new secret
	secretBytes := make([]byte, 32)
	_, err := rand.Read(secretBytes)
	if err != nil {
		log.Fatal("Failed to generate random bytes:", err)
	}

	newSecret := base64.URLEncoding.EncodeToString(secretBytes)
	fmt.Printf("Generated new JWT secret: %s\n", newSecret)

	// Read .env file
	envFile := ".env"
	content, err := ioutil.ReadFile(envFile)
	if err != nil {
		log.Fatal("Failed to read .env file:", err)
	}

	// Update JWT_SECRET line
	lines := strings.Split(string(content), "\n")
	updated := false

	for i, line := range lines {
		if strings.HasPrefix(line, "JWT_SECRET=") {
			lines[i] = fmt.Sprintf("JWT_SECRET=%s", newSecret)
			updated = true
			fmt.Println("✅ Updated JWT_SECRET in .env file")
			break
		}
	}

	if !updated {
		// Add JWT_SECRET if not found
		lines = append(lines, fmt.Sprintf("JWT_SECRET=%s", newSecret))
		fmt.Println("✅ Added JWT_SECRET to .env file")
	}

	// Write back to file
	newContent := strings.Join(lines, "\n")
	err = ioutil.WriteFile(envFile, []byte(newContent), 0644)
	if err != nil {
		log.Fatal("Failed to write .env file:", err)
	}

	fmt.Println("🎉 JWT secret successfully updated!")

	// Ask if user wants to restart server
	fmt.Print("\n⚠️  Remember to restart your server to use the new secret!\n")
	fmt.Print("🔄 Do you want to see the updated .env file? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		fmt.Println("\n📄 Updated .env file:")
		fmt.Println("=====================")
		updatedContent, _ := ioutil.ReadFile(envFile)
		fmt.Print(string(updatedContent))
	}
}
