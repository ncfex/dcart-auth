package user

import (
	"time"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

type UserRegisteredEventPayload struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func NewUserRegisteredEvent(aggregateID string, username, passwordHash string) shared.Event {
	return shared.BaseEvent{
		AggregateID:   aggregateID,
		AggregateType: "USER",
		EventType:     "USER_REGISTERED",
		Version:       1,
		Timestamp:     time.Now(),
		Payload:       UserRegisteredEventPayload{Username: username, PasswordHash: passwordHash},
	}
}

type UserPasswordChangedEventPayload struct {
	NewPasswordHash string `json:"new_password_hash"`
}

func NewUserPasswordChangedEvent(aggregateID string, newPasswordHash string, version int) shared.Event {
	return shared.BaseEvent{
		AggregateID:   aggregateID,
		AggregateType: "USER",
		EventType:     "USER_PASSWORD_CHANGED",
		Version:       version,
		Timestamp:     time.Now(),
		Payload:       UserPasswordChangedEventPayload{NewPasswordHash: newPasswordHash},
	}
}
