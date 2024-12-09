package user

import (
	"time"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

type UserRegisteredEvent struct {
	shared.BaseEvent
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func NewUserRegisteredEvent(aggregateID string, username, passwordHash string) *UserRegisteredEvent {
	return &UserRegisteredEvent{
		BaseEvent: shared.BaseEvent{
			AggregateID:   aggregateID,
			AggregateType: "USER",
			EventType:     string(EventTypeUserRegistered),
			Version:       1,
			Timestamp:     time.Now(),
		},
		Username:     username,
		PasswordHash: passwordHash,
	}
}

type UserPasswordChangedEvent struct {
	shared.BaseEvent
	NewPasswordHash string `json:"new_password_hash"`
}

func NewUserPasswordChangedEvent(aggregateID string, newPasswordHash string, version int) *UserPasswordChangedEvent {
	return &UserPasswordChangedEvent{
		BaseEvent: shared.BaseEvent{
			AggregateID:   aggregateID,
			AggregateType: "USER",
			EventType:     string(EventTypeUserPasswordChanged),
			Version:       version,
			Timestamp:     time.Now(),
		},
		NewPasswordHash: newPasswordHash,
	}
}
