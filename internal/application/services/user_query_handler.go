package services

import (
	"context"

	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	"github.com/ncfex/dcart-auth/internal/application/queries"
	"github.com/ncfex/dcart-auth/internal/domain/shared"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type userQueryHandler struct {
	eventStore outbound.EventStore
}

func NewUserQueryHandler(eventStore outbound.EventStore) *userQueryHandler {
	return &userQueryHandler{
		eventStore: eventStore,
	}
}

func (h *userQueryHandler) HandleGetUserById(ctx context.Context, query queries.GetUserByIDQuery) (*userDomain.User, error) {
	events, err := h.eventStore.GetEvents(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, userDomain.ErrUserNotFound
	}

	user := &userDomain.User{
		BaseAggregateRoot: shared.BaseAggregateRoot{
			Changes: []shared.Event{},
		},
	}
	for _, event := range events {
		user.Apply(event)
	}

	return user, nil
}

func (h *userQueryHandler) HandleGetUserByUsername(ctx context.Context, query queries.GetUserByUsernameQuery) (*userDomain.User, error) {
	events, err := h.eventStore.GetEventsByUsername(ctx, query.Username)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, userDomain.ErrUserNotFound
	}

	user := &userDomain.User{
		BaseAggregateRoot: shared.BaseAggregateRoot{
			Changes: []shared.Event{},
		},
	}
	for _, event := range events {
		user.Apply(event)
	}

	return user, nil
}
