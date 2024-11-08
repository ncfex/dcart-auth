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

type HexTokenService interface {
	Make() (string, error)
}

type service struct {
	prefix      string
	tokenLength int
}

func NewHexTokenService(prefix string, tokenLength int) HexTokenService {
	return &service{
		prefix:      prefix,
		tokenLength: tokenLength,
	}
}

func (s *service) Make() (string, error) {
	token := make([]byte, s.tokenLength)
	_, err := rand.Read(token)
	if err != nil {
		return "", ErrTokenGenerationFailed
	}

	tokenString := fmt.Sprintf("%s%s", s.prefix, hex.EncodeToString(token))
	return tokenString, nil
}
