package utils

import (
	"math/rand"
	"time"
)

func RandomString(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().Unix()))

	generated := make([]byte, length)

	for i := range generated {
		generated[i] = charset[seededRand.Intn(len(charset)-1)]
	}

	return string(generated)
}

func RandomAlphaNumericString(length int) string {
	charset := "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ0123456789"

	return RandomString(length, charset)
}

func RandomNumericToken() string {
	length := 6
	charset := "0123456789"

	return RandomString(length, charset)
}
