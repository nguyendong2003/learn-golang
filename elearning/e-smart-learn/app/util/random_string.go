package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func GenerateOAuthState() string {
	return GenerateRandomString(16)
}

func GenerateRandomNumber(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	for i := 0; i < length; i++ {
		b[i] = (b[i] % 10) + '0'
	}

	return string(b)
}
