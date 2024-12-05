package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/security"

	"github.com/ncfex/dcart-auth/pkg/httputil/request"
	"github.com/ncfex/dcart-auth/pkg/httputil/response"
	"github.com/ncfex/dcart-auth/pkg/middleware"
)

// todo improve
func RequireJWTAuth(
	tokenValidator security.TokenValidator,
	responder response.Responder,
) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
			defer cancel()

			accessToken, err := request.GetBearerToken(r.Header)
			if err != nil {
				responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: missing or invalid token", err)
				return
			}

			userID, err := tokenValidator.Validate(accessToken)
			if err != nil {
				responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: invalid token", err)
				return
			}

			ctx = context.WithValue(ctx, request.ContextUserKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRefreshToken(
	tokenRepo outbound.TokenRepository,
	responder response.Responder,
) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
			defer cancel()

			refreshToken, err := request.GetBearerToken(r.Header)
			if err != nil {
				responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: missing or invalid refresh token", err)
				return
			}

			token, err := tokenRepo.GetByToken(ctx, refreshToken)
			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded):
					responder.RespondWithError(w, http.StatusGatewayTimeout, "Request timeout", err)
				default:
					responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: invalid refresh token", err)
				}
				return
			}

			if err = token.IsValid(); err != nil {
				responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: invalid refresh token", err)
				return
			}

			ctx = context.WithValue(ctx, request.ContextUserKey, token.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
