package token

import (
	"fmt"
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TokenValid(token string) error {
	need := os.Getenv("ACCESS_TOKEN")
	if need == "" {
		return fmt.Errorf("ACCESS_TOKEN environment variable is not defined")
	}
	if token != need {
		return fmt.Errorf("Token is invalid")
	}
	return nil
}
