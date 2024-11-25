package services

import (
	"context"
	"time"

	"github.com/ncfex/dcart-auth/internal/application/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type service struct {
	userRepo        outbound.UserRepository
	tokenRepo       outbound.TokenRepository
	accessTokenGen  outbound.TokenGenerator
	refreshTokenGen outbound.RefreshTokenGenerator
	// svc
	userService inbound.UserSevice
}

func NewAuthService(
	userRepo outbound.UserRepository,
	tokenRepo outbound.TokenRepository,
	accessTokenGen outbound.TokenGenerator,
	refreshTokenGen outbound.RefreshTokenGenerator,
	// svc
	userService inbound.UserSevice,
) *service {
	return &service{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		accessTokenGen:  accessTokenGen,
		refreshTokenGen: refreshTokenGen,
		userService:     userService,
	}
}

func (s *service) Register(ctx context.Context, username, password string) (*userDomain.User, error) {
	_, err := s.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, userDomain.ErrUserAlreadyExists
	}
	if err != userDomain.ErrUserNotFound {
		return nil, err
	}

	user, err := s.userService.CreateUser(username, password)
	if err != nil {
		return nil, userDomain.ErrInvalidCredentials
	}

	_, err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, username, password string) (*tokenDomain.TokenPair, error) {
	if username == "" || password == "" {
		return nil, userDomain.ErrInvalidCredentials
	}
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, userDomain.ErrInvalidCredentials
	}

	err = s.userService.ComparePassword(user.PasswordHash, password)
	if err != nil {
		return nil, userDomain.ErrInvalidCredentials
	}

	accessTokenString, err := s.accessTokenGen.Generate(user.ID, time.Minute*15)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := s.refreshTokenGen.Generate()
	if err != nil {
		return nil, err
	}

	refreshToken, err := tokenDomain.NewRefreshToken(refreshTokenString, user.ID)
	if err != nil {
		return nil, err
	}

	err = s.tokenRepo.StoreToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &tokenDomain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s *service) Refresh(ctx context.Context, tokenString string) (*tokenDomain.TokenPair, error) {
	if tokenString == "" {
		return nil, tokenDomain.ErrTokenInvalid
	}
	refreshToken, err := s.tokenRepo.GetTokenByTokenString(ctx, tokenString)
	if err != nil {
		return nil, tokenDomain.ErrTokenInvalid
	}

	user, err := s.userRepo.GetUserByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.accessTokenGen.Generate(user.ID, time.Minute*15)
	if err != nil {
		return nil, err
	}

	return &tokenDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *service) Logout(ctx context.Context, tokenString string) error {
	token, err := s.tokenRepo.GetTokenByTokenString(ctx, tokenString)
	if err != nil {
		return tokenDomain.ErrTokenInvalid
	}

	token.Revoke()

	err = s.tokenRepo.Save(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Validate(ctx context.Context, token string) (*userDomain.User, error) {
	if token == "" {
		return nil, tokenDomain.ErrTokenInvalid
	}

	userID, err := s.accessTokenGen.Validate(token)
	if err != nil {
		return nil, tokenDomain.ErrTokenInvalid
	}

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, userDomain.ErrUserNotFound
	}

	return user, nil
}
