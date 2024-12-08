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
	switch event.GetEventType() {
	case "USER_REGISTERED":
		return p.projectUserRegistered(ctx, event)
	default:
		return fmt.Errorf("unsupported event type: %s", event.GetEventType())
	}
}

func (p *MongoProjector) projectUserRegistered(ctx context.Context, event shared.Event) error {
	payload, ok := event.GetPayload().(user.UserRegisteredEventPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for USER_REGISTERED event")
	}

	collection := p.db.Collection(p.collectionName)

	filter := bson.M{"_id": event.GetAggregateID()}
	update := bson.M{
		"$setOnInsert": bson.M{
			"_id":        event.GetAggregateID(),
			"username":   payload.Username,
			"created_at": event.GetTimestamp(),
			"updated_at": event.GetTimestamp(),
			"version":    1,
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to project USER_REGISTERED event: %w", err)
	}

	return nil
}
