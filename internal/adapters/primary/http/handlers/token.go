package handlers

import (
	"net/http"

	"github.com/ncfex/dcart-auth/pkg/httputil/request"
)

func (h *handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	tokenPair, err := h.authenticationService.Refresh(r.Context(), refreshToken)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		Token: string(tokenPair.AccessToken),
	})
}

func (h *handler) validateToken(w http.ResponseWriter, r *http.Request) {
	token, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	user, err := h.authenticationService.Validate(r.Context(), token)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"valid": true,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
