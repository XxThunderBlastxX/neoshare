package utils

import (
	"strings"
	"unicode"
)

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
