package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type userService struct {
	passwordHasher outbound.PasswordHasher
	userRepo       outbound.UserRepository
}

func NewUserService(
	passwordHasher outbound.PasswordHasher,
	userRepo outbound.UserRepository,
) inbound.UserSevice {
	return &userService{
		passwordHasher: passwordHasher,
		userRepo:       userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, password string) (*userDomain.User, error) {
	_, err := s.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, fmt.Errorf("get user username: %w", userDomain.ErrUserAlreadyExists)
	}
	if err != userDomain.ErrUserNotFound {
		return nil, fmt.Errorf("get user username: %w", err)
	}

	user, err := userDomain.NewUser(username, password)
	if err != nil {
		return nil, fmt.Errorf("new user: %w", err)
	}

	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("password hash: %w", err)
	}

	user.SetHashedPassword(hashedPassword)

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (s *userService) ValidateWithCreds(ctx context.Context, username, password string) (*userDomain.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user username: %w", err)
	}

	if err := s.comparePassword(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("compare password: %w", err)
	}
	return user, nil
}

func (s *userService) ValidateWithID(ctx context.Context, userID string) (*userDomain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user id: %w", err)
	}
	return user, nil
}

func (s *userService) comparePassword(hashedPassword, password string) error {
	if err := s.passwordHasher.Compare(hashedPassword, password); err != nil {
		return fmt.Errorf("compare: %w", err)
	}
	return nil
}
