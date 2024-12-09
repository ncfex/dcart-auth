package user

import (
	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

const (
	EventTypeUserRegistered      shared.EventType = "user.registered"
	EventTypeUserPasswordChanged shared.EventType = "user.passwordChanged"
)

func RegisterEvents(registry shared.EventRegistry) {
	registry.RegisterEvent(EventTypeUserRegistered, func() shared.Event {
		return &UserRegisteredEvent{}
	})
	registry.RegisterEvent(EventTypeUserPasswordChanged, func() shared.Event {
		return &UserPasswordChangedEvent{}
	})
}
