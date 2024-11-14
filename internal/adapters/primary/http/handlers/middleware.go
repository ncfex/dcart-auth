package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ncfex/dcart-auth/internal/core/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/core/ports/outbound"

	"github.com/ncfex/dcart-auth/pkg/httputil/request"
	"github.com/ncfex/dcart-auth/pkg/httputil/response"
	"github.com/ncfex/dcart-auth/pkg/middleware"
)

func RequireJWTAuth(
	tokenGenerator inbound.TokenGenerator,
	tokenRepo outbound.TokenRepository,
	userRepo outbound.UserRepository,
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

			userID, err := tokenGenerator.Validate(accessToken)
			if err != nil {
				responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: invalid token", err)
				return
			}

			user, err := userRepo.GetUserByID(ctx, userID)
			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded):
					responder.RespondWithError(w, http.StatusGatewayTimeout, "Request timeout", err)
				default:
					responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: user not found", err)
				}
				return
			}

			ctx = context.WithValue(ctx, request.ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRefreshToken(
	tokenRepo outbound.TokenRepository,
	userRepo outbound.UserRepository,
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

			user, err := tokenRepo.GetUserFromToken(ctx, refreshToken)
			if err != nil {
				switch {
				case errors.Is(err, context.DeadlineExceeded):
					responder.RespondWithError(w, http.StatusGatewayTimeout, "Request timeout", err)
				default:
					responder.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: invalid refresh token", err)
				}
				return
			}

			ctx = context.WithValue(ctx, request.ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
