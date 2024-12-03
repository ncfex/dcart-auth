package services

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/commands"
	"github.com/ncfex/dcart-auth/internal/application/ports/infra"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type userCommandHandler struct {
	eventStore infra.EventStore // separate read/write
}

func NewUserCommandHandler(eventStore infra.EventStore) *userCommandHandler {
	return &userCommandHandler{
		eventStore: eventStore,
	}
}

func (h *userCommandHandler) HandleRegisterUser(ctx context.Context, cmd commands.RegisterUserCommand) (*userDomain.User, error) {
	user, err := userDomain.NewUser(cmd.Username, cmd.Password)
	if err != nil {
		return nil, err
	}

	if err := h.eventStore.SaveEvents(ctx, user.GetID(), user.GetUncommittedChanges()); err != nil {
		return nil, err
	}

	user.ClearUncommittedChanges()
	return user, nil
}
