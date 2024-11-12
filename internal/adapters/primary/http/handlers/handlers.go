package handlers

import (
	"log"
	"net/http"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/middleware"
	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/response"
	"github.com/ncfex/dcart-auth/internal/ports"
)

type handler struct {
	logger            *log.Logger
	responder         response.Responder
	userAuthenticator ports.UserAuthenticator
	tokenManager      ports.TokenManager
	tokenRepo         ports.TokenRepository
	userRepo          ports.UserRepository
}

func NewHandler(
	logger *log.Logger,
	responder response.Responder,
	userAuthenticator ports.UserAuthenticator,
	tokenManager ports.TokenManager,
	tokenRepo ports.TokenRepository,
	userRepo ports.UserRepository,
) *handler {
	return &handler{
		logger:            logger,
		userAuthenticator: userAuthenticator,
		responder:         responder,
		tokenManager:      tokenManager,
		tokenRepo:         tokenRepo,
		userRepo:          userRepo,
	}
}

// TODO - add use or addRoute function
func (h *handler) Router() *http.ServeMux {
	mux := http.NewServeMux()

	loggingMiddleware := middleware.Logging(h.logger)
	recoveryMiddleware := middleware.Recovery(h.responder, h.logger)

	publicChain := middleware.Chain(
		loggingMiddleware,
		recoveryMiddleware,
	)

	refreshTokenRequiredChain := middleware.Chain(
		middleware.RequireRefreshToken(
			h.tokenManager,
			h.tokenRepo,
			h.userRepo,
			h.responder,
		),
		loggingMiddleware,
		recoveryMiddleware,
	)

	accessTokenProtectedChain := middleware.Chain(
		middleware.RequireJWTAuth(
			h.tokenManager,
			h.tokenRepo,
			h.userRepo,
			h.responder,
		),
		loggingMiddleware,
		recoveryMiddleware,
	)

	// public
	mux.Handle("POST /register", publicChain(http.HandlerFunc(h.register)))
	mux.Handle("POST /login", publicChain(http.HandlerFunc(h.login)))

	// protected
	mux.Handle("GET /profile", accessTokenProtectedChain(http.HandlerFunc(h.profile)))
	mux.Handle("POST /validate", accessTokenProtectedChain(http.HandlerFunc(h.validateToken)))

	// refresh required
	mux.Handle("POST /refresh", refreshTokenRequiredChain(http.HandlerFunc(h.refreshToken)))
	mux.Handle("POST /logout", refreshTokenRequiredChain(http.HandlerFunc(h.logout)))

	return mux
}
