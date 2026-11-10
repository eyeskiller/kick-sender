package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// GeneratePKCE generates a cryptographically secure PKCE code verifier and challenge
func GeneratePKCE() (verifier string, challenge string, err error) {
	// Generate 32 random bytes for verifier
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Base64 URL encode without padding
	verifier = base64.RawURLEncoding.EncodeToString(bytes)

	// Generate challenge using S256
	hash := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(hash[:])

	return verifier, challenge, nil
}

// GenerateState generates a random state parameter for OAuth
func GenerateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateState validates that the returned state matches the expected state
func ValidateState(expected, actual string) bool {
	if len(expected) == 0 || len(actual) == 0 {
		return false
	}
	// Simple constant-time comparison
	if len(expected) != len(actual) {
		return false
	}
	var result byte
	for i := 0; i < len(expected); i++ {
		result |= expected[i] ^ actual[i]
	}
	return result == 0
}
