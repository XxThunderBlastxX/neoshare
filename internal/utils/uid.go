package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// GenerateUID generates a UID using a hash function
func GenerateUID(input string) string {
	// Create a new SHA256 hash
	hash := sha256.New()
	// Write the input data to the hash
	hash.Write([]byte(input))
	// Get the hash sum
	sum := hash.Sum(nil)
	// Encode the hash sum to a base64 string
	uid := base64.RawURLEncoding.EncodeToString(sum)

	return uid[:8]
}
