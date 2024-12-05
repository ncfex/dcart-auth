package readmodel

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	"github.com/ncfex/dcart-auth/internal/application/queries"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserQueryHandler struct {
	db *mongo.Database
}

func (h *UserQueryHandler) GetUserByID(ctx context.Context, query queries.GetUserByIDQuery) (*types.UserResponse, error) {
	var userRM UserReadModel
	err := h.db.Collection("users").FindOne(ctx, bson.M{"_id": query.UserID}).Decode(&userRM)
	if err == mongo.ErrNoDocuments {
		return nil, userDomain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &types.UserResponse{
		ID:       userRM.ID,
		Username: userRM.Username,
	}, nil
}

func (h *UserQueryHandler) GetUserByUsername(ctx context.Context, query queries.GetUserByUsernameQuery) (*types.UserResponse, error) {
	var userRM UserReadModel
	err := h.db.Collection("users").FindOne(ctx, bson.M{"username": query.Username}).Decode(&userRM)
	if err == mongo.ErrNoDocuments {
		return nil, userDomain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &types.UserResponse{
		ID:       userRM.ID,
		Username: userRM.Username,
	}, nil
}