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
) inbound.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) CreateUser(ctx context.Context, req inbound.RegisterRequest) (*inbound.UserDTO, error) {
	if _, err := us.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, fmt.Errorf("get user username: %w", userDomain.ErrUserAlreadyExists)
	} else if err != userDomain.ErrUserNotFound {
		return nil, fmt.Errorf("get user username: %w", err)
	}

	user, err := userDomain.NewUser(req.Username, req.Password)
	if err != nil {
		return nil, fmt.Errorf("new user: %w", err)
	}

	if err := us.userRepo.Add(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &inbound.UserDTO{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (us *userService) VerifyWithCreds(ctx context.Context, req inbound.LoginRequest) (*inbound.UserDTO, error) {
	user, err := us.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("get user username: %w", err)
	}

	if ok := user.Authenticate(req.Password); !ok {
		return nil, fmt.Errorf("authenticate wrong password")
	}
	return &inbound.UserDTO{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (us *userService) GetExistingUser(ctx context.Context, userID string) (*inbound.UserDTO, error) {
	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user id: %w", err)
	}
	return &inbound.UserDTO{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
