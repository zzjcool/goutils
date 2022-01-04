package goutils

import (
	"crypto/rand"
	"fmt"
)

func GetRandomString(n uint) string {
	randBytes := make([]byte, n/2)
	_, _ = rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
