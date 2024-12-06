package command

import (
	"context"
	"fmt"
	"log"

	"github.com/ncfex/dcart-auth/internal/application/ports/id"
	"github.com/ncfex/dcart-auth/internal/application/ports/primary/command"
	"github.com/ncfex/dcart-auth/internal/application/ports/secondary"
	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type UserCommandHandler struct {
	eventStore     secondary.EventStore
	eventPublisher secondary.EventPublisher
	idGenerator    id.IDGenerator
}

func NewUserCommandHandler(
	eventStore secondary.EventStore,
	eventPublisher secondary.EventPublisher,
	idGenerator id.IDGenerator,
) command.UserCommandPort {
	return &UserCommandHandler{
		eventStore:     eventStore,
		eventPublisher: eventPublisher,
		idGenerator:    idGenerator,
	}
}

func (h *UserCommandHandler) RegisterUser(ctx context.Context, cmd command.RegisterUserCommand) (*types.UserResponse, error) {
	userID := h.idGenerator.GenerateFromData([]byte(cmd.Username))

	events, err := h.eventStore.GetEvents(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}
	if len(events) > 0 {
		return nil, userDomain.ErrUserAlreadyExists
	}

	newUser, err := userDomain.NewUser(userID, cmd.Username, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	events = newUser.GetUncommittedChanges()
	if err := h.eventStore.SaveEvents(ctx, userID, events); err != nil {
		return nil, fmt.Errorf("saving events: %w", err)
	}

	for _, event := range events {
		if err := h.eventPublisher.PublishEvent(ctx, event); err != nil {
			log.Printf("error publishing event: %v", err)
		}
	}

	return &types.UserResponse{
		ID:       userID,
		Username: cmd.Username,
	}, nil
}

func (h *UserCommandHandler) AuthenticateUser(ctx context.Context, cmd command.AuthenticateUserCommand) (*types.UserResponse, error) {
	userID := h.idGenerator.GenerateFromData([]byte(cmd.Username))

	events, err := h.eventStore.GetEvents(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}

	currentUser, err := userDomain.ReconstructFromEvents(events)
	if err != nil {
		return nil, fmt.Errorf("applying events: %w", err)
	}

	if ok := currentUser.Authenticate(cmd.Password); !ok {
		return nil, fmt.Errorf("saving events: %w", err)
	}

	// maybe publish event

	return &types.UserResponse{
		ID:       currentUser.ID,
		Username: currentUser.Username,
	}, nil
}

func (h *UserCommandHandler) ChangePassword(ctx context.Context, cmd command.ChangePasswordCommand) error {
	events, err := h.eventStore.GetEvents(ctx, cmd.UserID)
	if err != nil {
		return fmt.Errorf("loading events: %w", err)
	}

	currentUser, err := userDomain.ReconstructFromEvents(events)
	if err != nil {
		return fmt.Errorf("applying events: %w", err)
	}

	if err := currentUser.ChangePassword(cmd.OldPassword, cmd.NewPassword); err != nil {
		return fmt.Errorf("changing password: %w", err)
	}

	newEvents := currentUser.GetUncommittedChanges()
	if err := h.eventStore.SaveEvents(ctx, cmd.UserID, newEvents); err != nil {
		return fmt.Errorf("saving events: %w", err)
	}

	for _, event := range newEvents {
		if err := h.eventPublisher.PublishEvent(ctx, event); err != nil {
			log.Printf("error publishing event: %v", err)
		}
	}

	return nil
}
