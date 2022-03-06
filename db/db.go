package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Queries interface {
	CreateTokens(uid string, duration time.Duration) (map[string]string, error)

	Find(ctx context.Context, id string) (user *User, err error)
	Replace(ctx context.Context, user *User, arg interface{}) (*mongo.UpdateResult, error)
}
