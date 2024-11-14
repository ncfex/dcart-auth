package refresh

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrTokenGenerationFailed = errors.New("token generation failed")
)

type hexRefreshGenerator struct {
	prefix      string
	tokenLength int
}

func NewHexRefreshGenerator(prefix string, tokenLength int) *hexRefreshGenerator {
	return &hexRefreshGenerator{
		prefix:      prefix,
		tokenLength: tokenLength,
	}
}

func (g *hexRefreshGenerator) Generate() (string, error) {
	token := make([]byte, g.tokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return "", ErrTokenGenerationFailed
	}

	tokenString := fmt.Sprintf("%s%s", g.prefix, hex.EncodeToString(token))
	return tokenString, nil
}
