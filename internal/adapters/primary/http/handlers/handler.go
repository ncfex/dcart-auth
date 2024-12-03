package handlers

import (
	"log"
	"net/http"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/middlewares"
	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/infra"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"

	"github.com/ncfex/dcart-auth/pkg/httputil/response"
	"github.com/ncfex/dcart-auth/pkg/middleware"
)

type handler struct {
	logger                *log.Logger
	responder             response.Responder
	authenticationService inbound.AuthenticationService
	tokenManager          outbound.TokenGeneratorValidator
	tokenRepo             infra.TokenRepository
	eventStore            infra.EventStore
}

func NewHandler(
	logger *log.Logger,
	responder response.Responder,
	authenticationService inbound.AuthenticationService,
	tokenManager outbound.TokenGeneratorValidator,
	tokenRepo infra.TokenRepository,
	eventStore infra.EventStore,
) *handler {
	return &handler{
		logger:                logger,
		authenticationService: authenticationService,
		responder:             responder,
		tokenManager:          tokenManager,
		tokenRepo:             tokenRepo,
		eventStore:            eventStore,
	}
}

func (h *handler) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	loggingMiddleware := middleware.Logging(h.logger)
	recoveryMiddleware := middleware.Recovery(h.responder, h.logger)

	publicChain := middleware.Chain(
		loggingMiddleware,
		recoveryMiddleware,
	)

	refreshTokenRequiredChain := middleware.Chain(
		middlewares.RequireRefreshToken(
			h.tokenRepo,
			h.responder,
		),
		loggingMiddleware,
		recoveryMiddleware,
	)

	accessTokenProtectedChain := middleware.Chain(
		middlewares.RequireJWTAuth(
			h.tokenManager,
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
