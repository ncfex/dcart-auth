package handlers

import (
	"net/http"

	"github.com/ncfex/dcart-auth/internal/adapters/primary/http/request"
	"github.com/ncfex/dcart-auth/internal/domain"
)

func (h *handler) profile(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User domain.User `json:"user"`
	}

	user, exists := request.GetUserFromContext(r.Context())
	if !exists {
		h.responder.RespondWithError(w, http.StatusNotFound, domain.ErrUserNotFound.Error(), domain.ErrUserNotFound)
		return
	}

	h.responder.RespondWithJSON(w, http.StatusOK, response{
		User: domain.User{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}
