package request

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrRequestFailed        = errors.New("request failed")
	ErrTimeout              = errors.New("timeout")
	ErrNoAuthHeaderIncluded = errors.New("no auth header included")
	ErrMalformedAuthHeader  = errors.New("malformed authorization header")
)

const (
	BearerPrefix        string = "Bearer "
	AuthorizationHeader string = "Authorization"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get(AuthorizationHeader)
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", ErrMalformedAuthHeader
	}

	token := authHeader[len(BearerPrefix):]
	if token == "" {
		return "", ErrMalformedAuthHeader
	}

	return token, nil
}
