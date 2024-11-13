package handlers

import (
	"net/http"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/request"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
)

func (h *handler) profile(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User userDomain.User `json:"user"`
	}

	user, exists := request.GetUserFromContext(r.Context())
	if !exists {
		h.responder.RespondWithError(w, http.StatusNotFound, userDomain.ErrUserNotFound.Error(), userDomain.ErrUserNotFound)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		User: userDomain.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}
