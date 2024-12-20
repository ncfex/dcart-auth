package handlers

import (
	"net/http"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	"github.com/ncfex/dcart-auth/pkg/httputil/request"
)

func (h *handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	tokenPairResponse, err := h.authenticationService.Refresh(r.Context(), types.TokenRequest{
		Token: refreshToken,
	})
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, tokenPairResponse)
}

func (h *handler) validateToken(w http.ResponseWriter, r *http.Request) {
	accessToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	validateResponse, err := h.authenticationService.Validate(r.Context(), types.TokenRequest{
		Token: accessToken,
	})
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, validateResponse)
}
