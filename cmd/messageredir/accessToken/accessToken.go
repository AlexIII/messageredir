package accessToken

import (
	"crypto/rand"
	"encoding/base64"
)

// Generates a secure random access token
func Generate(length int) (string, error) {
	// Create a byte slice to hold the random bytes
	bytes := make([]byte, length)

	// Fill the byte slice with random data
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a Base64 string
	token := base64.RawURLEncoding.EncodeToString(bytes)

	return token, nil
}
