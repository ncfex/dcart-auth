package services

import (
	"github.com/google/uuid"
	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type userService struct {
	passwordHasher outbound.PasswordHasher
}

func NewUserService(
	passwordHasher outbound.PasswordHasher,
) inbound.UserSevice {
	return &userService{
		passwordHasher: passwordHasher,
	}
}

func (s *userService) CreateUser(username, password string) (*userDomain.User, error) {
	newID := uuid.New().String()
	user, err := userDomain.NewUser(newID, username, password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}
	user.SetHashedPassword(hashedPassword)

	return user, nil
}

func (s *userService) ComparePassword(hashedPassword, password string) error {
	err := s.passwordHasher.Compare(hashedPassword, password)
	if err != nil {
		return err
	}

	return nil
}
