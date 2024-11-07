package refresh

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
		return "", err
	}
	return fmt.Sprintf("%s%s", s.prefix, hex.EncodeToString(token)), nil
}
