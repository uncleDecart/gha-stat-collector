package token_test

import (
	"regexp"
	"testing"

	"github.com/uncleDecart/gha-stat-collector/pkg/token"
	"gotest.tools/assert"
)

func TestGenerateToken(t *testing.T) {
	tokenLength := 64
	got := token.GenerateToken(tokenLength)
	assert.Equal(t, tokenLength, len(got), "Generated token should be specified length")

	symbol_mask := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(got)

	assert.Equal(t, true, symbol_mask, "Generated token should consist only of symbols specified")
}

func TestTokenValid(t *testing.T) {
	tokenValue := "qwerty"
	t.Setenv("ACCESS_TOKEN", tokenValue)
	err := token.TokenValid(tokenValue)
	assert.Equal(t, nil, err, "Token should be valid")
}

func BenchmarkGenerateToken(b *testing.B) {
	tokenLength := 64
	for i := 0; i < b.N; i++ {
		token.GenerateToken(tokenLength)
	}
}
