package authentication

import (
	"context"
	"fmt"
	"time"

	tokenDomain "github.com/ncfex/dcart-auth/internal/core/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/core/domain/user"
	"github.com/ncfex/dcart-auth/internal/core/ports/inbound"
	"github.com/ncfex/dcart-auth/internal/core/ports/outbound"
	"github.com/ncfex/dcart-auth/pkg/services/auth/tokens/refresh"
)

type service struct {
	userRepo            outbound.UserRepository
	tokenRepo           outbound.TokenRepository
	passwordEncrypter   inbound.PasswordEncrypter
	accessTokenManager  inbound.TokenManager
	refreshTokenManager refresh.HexTokenService
}

func NewAuthService(
	userRepo outbound.UserRepository,
	tokenRepo outbound.TokenRepository,
	passwordEncrypter inbound.PasswordEncrypter,
	accessTokenManager inbound.TokenManager,
	refreshTokenManager refresh.HexTokenService,
) inbound.UserAuthenticator {
	return &service{
		userRepo:            userRepo,
		tokenRepo:           tokenRepo,
		passwordEncrypter:   passwordEncrypter,
		accessTokenManager:  accessTokenManager,
		refreshTokenManager: refreshTokenManager,
	}
}

func (s *service) Register(ctx context.Context, username, password string) (*userDomain.User, error) {
	if username == "" || password == "" {
		return nil, userDomain.ErrInvalidCredentials
	}

	_, err := s.userRepo.GetUserByUsername(ctx, username)
	if err == nil {
		return nil, userDomain.ErrUserAlreadyExists
	}
	if err != userDomain.ErrUserNotFound {
		return nil, err
	}

	hashedPassword, err := s.passwordEncrypter.Hash(password)
	if err != nil {
		return nil, err
	}

	newUser := &userDomain.User{
		Username:     username,
		PasswordHash: hashedPassword,
	}

	return s.userRepo.CreateUser(ctx, newUser)
}

func (s *service) Login(ctx context.Context, username, password string) (*tokenDomain.TokenPair, error) {
	if username == "" || password == "" {
		return nil, userDomain.ErrInvalidCredentials
	}
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, userDomain.ErrInvalidCredentials
	}

	err = s.passwordEncrypter.Compare(user.PasswordHash, password)
	if err != nil {
		return nil, userDomain.ErrInvalidCredentials
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

	return &tokenDomain.TokenPair{
		AccessToken:  tokenDomain.Token(accessToken),
		RefreshToken: tokenDomain.Token(refreshToken),
	}, nil
}

func (s *service) Refresh(ctx context.Context, token string) (*tokenDomain.TokenPair, error) {
	if token == "" {
		return nil, tokenDomain.ErrTokenInvalid
	}

	ref, err := s.tokenRepo.GetTokenByTokenString(ctx, token)
	if err != nil {
		return nil, tokenDomain.ErrTokenInvalid
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

	return &tokenDomain.TokenPair{
		AccessToken:  tokenDomain.Token(accessToken),
		RefreshToken: tokenDomain.Token(token),
	}, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return tokenDomain.ErrTokenInvalid
	}

	_, err := s.tokenRepo.GetTokenByTokenString(ctx, token)
	if err != nil {
		return tokenDomain.ErrTokenInvalid
	}

	err = s.tokenRepo.RevokeToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
