package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	"github.com/ncfex/dcart-auth/pkg/httputil/request"
)

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
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
	var req types.LoginRequest
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

func (h *handler) changePassword(w http.ResponseWriter, r *http.Request) {
	var req types.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if err := h.authenticationService.ChangePassword(r.Context(), req); err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	err = h.authenticationService.Logout(r.Context(), types.TokenRequest{
		Token: refreshToken,
	})
	if err != nil {
		h.responder.RespondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
