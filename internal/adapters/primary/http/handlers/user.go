package handlers

import (
	"net/http"

	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
	"github.com/ncfex/dcart-auth/pkg/httputil/request"
)

// todo improve
func (h *handler) profile(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User userDomain.User `json:"user"`
	}

	userID, exists := request.GetDataFromContext[string](r.Context(), request.ContextUserKey)
	if !exists {
		h.responder.RespondWithError(w, http.StatusNotFound, userDomain.ErrUserNotFound.Error(), userDomain.ErrUserNotFound)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		User: userDomain.User{
			ID: *userID,
		},
	})
}
