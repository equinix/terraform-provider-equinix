package config

import (
	"math/rand"
	"time"
)

const allowed_charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$&@"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func correlationIdWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func correlationId(length int) string {
	return correlationIdWithCharset(length, allowed_charset)
}
