package middleware

import (
	"fmt"
	"go-get-backend/models"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Encoder creates JWT token from user data
func Encoder(userData models.JWTPayload) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-default-secret-key" // Default untuk development
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"id":    userData.ID,
		"email": userData.Email,
		"role":  userData.Role,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 7 days
		"iat":   time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Decoder verifies and decodes JWT token
func Decoder(tokenString string) (*models.JWTPayload, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "your-default-secret-key" // Default untuk development
	}

	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userData := &models.JWTPayload{
			ID:    claims["id"].(string),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}
		return userData, nil
	}

	return nil, fmt.Errorf("invalid token")
}
