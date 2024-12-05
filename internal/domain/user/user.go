package user

import (
	"errors"
	"time"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
)

var (
	ErrInvalidUser        = errors.New("invalid user")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type User struct {
	shared.BaseAggregateRoot
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(userID, username, rawPassword string) (*User, error) {
	if err := validateUserName(username); err != nil {
		return nil, err
	}

	password, err := NewPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := password.Hash()
	if err != nil {
		return nil, err
	}

	user := &User{
		BaseAggregateRoot: shared.BaseAggregateRoot{
			ID:      userID,
			Version: 0,
		},
	}

	event := NewUserRegisteredEvent(user.ID, username, hashedPassword)
	user.Apply(event)
	user.Changes = append(user.Changes, event)

	return user, nil
}

func (u *User) Authenticate(rawPassword string) bool {
	password := Password(rawPassword)
	return password.Matches(u.PasswordHash)
}

func (u *User) ChangePassword(rawOldPassword, rawNewPassword string) error {
	oldPwd := Password(rawOldPassword)
	if !oldPwd.Matches(u.PasswordHash) {
		return ErrInvalidPassword
	}

	newPassword, err := NewPassword(rawNewPassword)
	if err != nil {
		return err
	}

	newPasswordHash, err := newPassword.Hash()
	if err != nil {
		return err
	}

	u.PasswordHash = newPasswordHash
	event := NewUserPasswordChangedEvent(u.ID, u.PasswordHash, u.GetVersion())
	u.Apply(event)
	u.Changes = append(u.Changes, event)

	return nil
}

// todo - use value object
func validateUserName(username string) error {
	if username == "" {
		return ErrInvalidCredentials
	}
	return nil
}

func (u *User) Apply(event shared.Event) {
	u.ID = event.GetAggregateID()
	u.Version = event.GetVersion()

	switch event.GetEventType() {
	case "USER_REGISTERED":
		payload := event.GetPayload().(UserRegisteredEventPayload)
		u.Username = payload.Username
		u.PasswordHash = payload.PasswordHash
		u.CreatedAt = event.GetTimestamp()
		u.UpdatedAt = event.GetTimestamp()
	case "USER_PASSWORD_CHANGED":
		payload := event.GetPayload().(UserPasswordChangedEventPayload)
		u.PasswordHash = payload.NewPasswordHash
		u.UpdatedAt = event.GetTimestamp()
	}
}
