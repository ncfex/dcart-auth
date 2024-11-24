package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ncfex/dcart-auth/internal/application/ports/outbound"
	tokenDomain "github.com/ncfex/dcart-auth/internal/domain/token"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"
)

type service struct {
	userRepo        outbound.UserRepository
	tokenRepo       outbound.TokenRepository
	passwordHasher  outbound.PasswordHasher
	accessTokenGen  outbound.TokenGenerator
	refreshTokenGen outbound.RefreshTokenGenerator
}

func NewAuthService(
	userRepo outbound.UserRepository,
	tokenRepo outbound.TokenRepository,
	passwordHasher outbound.PasswordHasher,
	accessTokenGen outbound.TokenGenerator,
	refreshTokenGen outbound.RefreshTokenGenerator,
) *service {
	return &service{
		userRepo:        userRepo,
		tokenRepo:       tokenRepo,
		passwordHasher:  passwordHasher,
		accessTokenGen:  accessTokenGen,
		refreshTokenGen: refreshTokenGen,
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

	user, err := userDomain.New(username, password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}
	user.SetHashedPassword(hashedPassword)

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

	err = s.passwordHasher.Compare(user.PasswordHash, password)
	if err != nil {
		return nil, userDomain.ErrInvalidCredentials
	}

	accessToken, err := s.accessTokenGen.Generate(user.ID, time.Minute*15)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshTokenGen.Generate()
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

	accessToken, err := s.accessTokenGen.Generate(user.ID, time.Minute*15)
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
