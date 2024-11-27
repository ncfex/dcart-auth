package services

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type userService struct {
	userRepo outbound.UserRepository
}

func NewUserService(
	userRepo outbound.UserRepository,
) inbound.UserSevice {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) CreateUser(ctx context.Context, username, password string) (*userDomain.User, error) {
	_, err := us.userRepo.GetByUsername(ctx, username)
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

	if err := us.userRepo.Add(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (us *userService) ValidateWithCreds(ctx context.Context, username, password string) (*userDomain.User, error) {
	user, err := us.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user username: %w", err)
	}

	if ok := user.Authenticate(password); !ok {
		return nil, fmt.Errorf("authenticate wrong password")
	}
	return user, nil
}

func (us *userService) ValidateWithID(ctx context.Context, userID string) (*userDomain.User, error) {
	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user id: %w", err)
	}
	return user, nil
}
