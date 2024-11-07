package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/response"
)

var ErrInternalServerErrorStr = "internal server error"
var ErrInternalServerError = errors.New(ErrInternalServerErrorStr)

func Recovery(responder response.Responder, logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("Recovered from panic: %v", err)
					responder.RespondWithError(
						w,
						http.StatusInternalServerError,
						ErrInternalServerErrorStr,
						ErrInternalServerError,
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
