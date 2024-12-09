package mongodb

import (
	"context"
	"fmt"

	"github.com/ncfex/dcart-auth/internal/domain/shared"
	"github.com/ncfex/dcart-auth/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoProjector struct {
	db             *mongo.Database
	collectionName string
}

func NewMongoProjector(db *mongo.Database, collectionName string) *MongoProjector {
	return &MongoProjector{
		db:             db,
		collectionName: collectionName,
	}
}

func (p *MongoProjector) ProjectEvent(ctx context.Context, event shared.Event) error {
	switch e := event.(type) {
	case *user.UserRegisteredEvent:
		return p.projectUserRegistered(ctx, e)
	case *user.UserPasswordChangedEvent:
		return p.projectUserPasswordChanged(ctx, e)
	default:
		return fmt.Errorf("unsupported event type: %s", event.GetEventType())
	}
}

func (p *MongoProjector) projectUserRegistered(ctx context.Context, event *user.UserRegisteredEvent) error {
	collection := p.db.Collection(p.collectionName)

	filter := bson.M{"_id": event.GetAggregateID()}
	update := bson.M{
		"$setOnInsert": bson.M{
			"_id":           event.GetAggregateID(),
			"username":      event.Username,
			"password_hash": event.PasswordHash,
			"created_at":    event.GetTimestamp(),
			"updated_at":    event.GetTimestamp(),
			"version":       1,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (p *MongoProjector) projectUserPasswordChanged(ctx context.Context, event *user.UserPasswordChangedEvent) error {
	collection := p.db.Collection(p.collectionName)

	filter := bson.M{"_id": event.GetAggregateID()}
	update := bson.M{
		"$set": bson.M{
			"password_hash": event.NewPasswordHash,
			"updated_at":    event.GetTimestamp(),
			"version":       event.GetVersion(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}
