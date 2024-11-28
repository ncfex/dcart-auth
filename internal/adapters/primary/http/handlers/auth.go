package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/pkg/httputil/request"
)

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var req inbound.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	userResponse, err := h.authenticationService.Register(r.Context(), req)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusCreated, userResponse)
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var req inbound.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	tokenPairResponse, err := h.authenticationService.Login(r.Context(), req)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, tokenPairResponse)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	err = h.authenticationService.Logout(r.Context(), inbound.LogoutRequest{
		TokenString: refreshToken,
	})
	if err != nil {
		h.responder.RespondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
