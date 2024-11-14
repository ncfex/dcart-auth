package handlers

import (
	"log"
	"net/http"

	"github.com/ncfex/dcart-auth/internal/core/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/core/ports/outbound"

	"github.com/ncfex/dcart-auth/pkg/httputil/response"
	"github.com/ncfex/dcart-auth/pkg/middleware"
)

type handler struct {
	logger                *log.Logger
	responder             response.Responder
	authenticationService inbound.AuthenticationService
	tokenGenerator        inbound.TokenGenerator
	tokenRepo             outbound.TokenRepository
	userRepo              outbound.UserRepository
}

func NewHandler(
	logger *log.Logger,
	responder response.Responder,
	authenticationService inbound.AuthenticationService,
	tokenGenerator inbound.TokenGenerator,
	tokenRepo outbound.TokenRepository,
	userRepo outbound.UserRepository,
) *handler {
	return &handler{
		logger:                logger,
		authenticationService: authenticationService,
		responder:             responder,
		tokenGenerator:        tokenGenerator,
		tokenRepo:             tokenRepo,
		userRepo:              userRepo,
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
		RequireRefreshToken(
			h.tokenRepo,
			h.userRepo,
			h.responder,
		),
		loggingMiddleware,
		recoveryMiddleware,
	)

	accessTokenProtectedChain := middleware.Chain(
		RequireJWTAuth(
			h.tokenGenerator,
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
