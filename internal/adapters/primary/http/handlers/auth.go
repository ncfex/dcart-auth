package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/request"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
)

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		userDomain.User
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	createdUser, err := h.userAuthenticator.Register(r.Context(), params.Username, params.Password)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusCreated, response{
		User: userDomain.User{
			ID:           createdUser.ID,
			Username:     createdUser.Username,
			PasswordHash: createdUser.PasswordHash,
			CreatedAt:    createdUser.CreatedAt,
			UpdatedAt:    createdUser.UpdatedAt,
		},
	})
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.responder.RespondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	tokenPair, err := h.userAuthenticator.Login(r.Context(), params.Username, params.Password)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		Token:        string(tokenPair.AccessToken),
		RefreshToken: string(tokenPair.RefreshToken),
	})
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := request.GetBearerToken(r.Header)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	err = h.userAuthenticator.Logout(r.Context(), refreshToken)
	if err != nil {
		h.responder.RespondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
