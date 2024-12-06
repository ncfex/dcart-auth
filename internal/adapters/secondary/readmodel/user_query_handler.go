package readmodel

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/application/ports/primary/query"
	"github.com/ncfex/dcart-auth/internal/application/ports/types"
	userDomain "github.com/ncfex/dcart-auth/internal/domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserQueryHandler struct {
	db         *mongo.Database
	collection string
}

func NewUserQueryHandler(db *mongo.Database) *UserQueryHandler {
	return &UserQueryHandler{
		db:         db,
		collection: "users",
	}
}

func (h *UserQueryHandler) GetUserByID(ctx context.Context, query query.GetUserByIDQuery) (*types.UserResponse, error) {
	var userRM UserReadModel
	err := h.db.Collection(h.collection).FindOne(ctx, bson.M{"_id": query.UserID}).Decode(&userRM)
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

func (h *UserQueryHandler) GetUserByUsername(ctx context.Context, query query.GetUserByUsernameQuery) (*types.UserResponse, error) {
	var userRM UserReadModel
	err := h.db.Collection(h.collection).FindOne(ctx, bson.M{"username": query.Username}).Decode(&userRM)
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
