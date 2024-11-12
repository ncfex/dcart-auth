package authentication

import (
	"context"
	"fmt"
	"time"

	"github.com/ncfex/dcart-auth/internal/core/domain"
	"github.com/ncfex/dcart-auth/internal/core/ports"
	"github.com/ncfex/dcart-auth/internal/core/services/refresh"
)

type service struct {
	userRepo            ports.UserRepository
	tokenRepo           ports.TokenRepository
	passwordEncrypter   ports.PasswordEncrypter
	accessTokenManager  ports.TokenManager
	refreshTokenManager refresh.HexTokenService
}

func NewAuthService(
	userRepo ports.UserRepository,
	tokenRepo ports.TokenRepository,
	passwordEncrypter ports.PasswordEncrypter,
	accessTokenManager ports.TokenManager,
	refreshTokenManager refresh.HexTokenService,
) ports.UserAuthenticator {
	return &service{
		userRepo:            userRepo,
		tokenRepo:           tokenRepo,
		passwordEncrypter:   passwordEncrypter,
		accessTokenManager:  accessTokenManager,
		refreshTokenManager: refreshTokenManager,
	}
}

func (s *service) Register(ctx context.Context, username, password string) (*domain.User, error) {
	if username == "" || password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	_, err := s.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != domain.ErrUserNotFound {
		return nil, err
	}

	hashedPassword, err := s.passwordEncrypter.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *service) Login(ctx context.Context, username, password string) (*domain.TokenPair, error) {
	if username == "" || password == "" {
		return nil, domain.ErrInvalidCredentials
	}
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	err = s.passwordEncrypter.Compare(user.PasswordHash, password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.accessTokenManager.Make(&user.ID, time.Minute*15)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenManager.Make()
	if err != nil {
		return nil, err
	}

	err = s.tokenRepo.StoreToken(ctx, user, refreshToken)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  domain.Token(accessToken),
		RefreshToken: domain.Token(refreshToken),
	}, nil
}

func (s *service) Refresh(ctx context.Context, token string) (*domain.TokenPair, error) {
	if token == "" {
		return nil, domain.ErrTokenInvalid
	}

	ref, err := s.tokenRepo.GetTokenByTokenString(ctx, token)
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	fmt.Println(ref)

	user, err := s.tokenRepo.GetUserFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.accessTokenManager.Make(&user.ID, time.Minute*15)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  domain.Token(accessToken),
		RefreshToken: domain.Token(token),
	}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return domain.ErrTokenInvalid
	}

	_, err := s.tokenRepo.GetTokenByTokenString(ctx, token)
	if err != nil {
		return domain.ErrTokenInvalid
	}

	err = s.tokenRepo.RevokeToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
