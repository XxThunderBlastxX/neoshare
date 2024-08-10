package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"unicode"
)

// GenerateRandomState generates a random state string
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

// RemoveNonAsciiValue removes non-ascii or non-printable characters from a string
func RemoveNonAsciiValue(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)

	return s
}
