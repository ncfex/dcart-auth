package request

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
	ErrMalformedAuthHeader  = errors.New("malformed authorization header")
)

const BearerSchema = "Bearer"

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Fields(authHeader)
	if len(splitAuth) != 2 || splitAuth[0] != BearerSchema {
		return "", ErrMalformedAuthHeader
	}

	token := splitAuth[1]
	if token == "" {
		return "", ErrMalformedAuthHeader
	}

	return token, nil
}
