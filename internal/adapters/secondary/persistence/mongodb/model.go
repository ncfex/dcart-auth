package mongodb

import (
	"time"
)

type UserReadModel struct {
	ID        string    `bson:"_id"`
	Username  string    `bson:"username"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
