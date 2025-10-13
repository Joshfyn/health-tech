package tokenmaker

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExpiry = 1 * time.Hour
)

type TokenRequirements struct {
	Email  string
	UserID string
	Role   string
}

type TokenClaims struct {
	TokenRequirements
	jwt.RegisteredClaims
}

var jwtSecret []byte = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken creates a JWT with a 1-hour expiration time.
func GenerateToken(tokenRequirements TokenRequirements) (string, error) {
	// Token expires 1 hour from now
	expirationTime := time.Now().Add(tokenExpiry)

	// Create the claims
	claims := &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Expiration Time
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			// Issued At
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Subject (useful for identifying the purpose)
			Subject: "email_verification",
		},
	}
	if tokenRequirements.Email != "" {
		claims.TokenRequirements.Email = tokenRequirements.Email
	}
	if tokenRequirements.UserID != "" {
		claims.TokenRequirements.UserID = tokenRequirements.UserID
	}
	if tokenRequirements.Role != "" {
		claims.TokenRequirements.Role = tokenRequirements.Role
	}

	// Declare the token with the claims and the signing method (HMAC SHA256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing JWT token: %v", err)
		return "", err
	}

	return tokenString, nil
}

// ValidateVerificationToken validates the token string and returns the email if valid.
func ValidateVerificationToken(tokenString string) (TokenClaims, error) {
	claims := &TokenClaims{}

	// Parse the token, automatically checking the signature and expiration time (TTL)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify that the signing method is HMAC (required by the signature validation)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for validation
		return jwtSecret, nil
	})

	if err != nil {
		// If the token is invalid (bad signature) OR expired (TTL check failed), an error is returned here.
		// Common errors include: "token is expired" or "signature is invalid".
		return TokenClaims{}, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return TokenClaims{}, fmt.Errorf("token is invalid")
	}

	// TTL check is handled automatically by jwt.ParseWithClaims()
	return *claims, nil
}
