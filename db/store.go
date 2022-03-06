package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/iKayrat/auth-jwt-mongo/token"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	UserCollection *mongo.Collection
	Token          token.Maker
	Queries
}

// NewStore connects to database and
func NewStore(client *mongo.Client) Store {
	usersCollection := client.Database("test").Collection("users")

	jwttoken, err := token.NewJWTMaker()
	if err != nil {
		log.Fatal(err)
	}

	return Store{
		UserCollection: usersCollection,
		Token:          jwttoken,
	}
}

type User struct {
	UserID       string    `bson:"user_id"`
	RefreshToken string    `bson:"refresh_token"`
	RefreshExp   time.Time `bson:"refresh_exp"`
}

func (s *Store) FindID(ctx context.Context, id string) (user *User, err error) {
	err = s.UserCollection.FindOne(ctx, bson.M{"user_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return
}

func (s *Store) FindUser(ctx context.Context, key, value string) (user *User, err error) {
	err = s.UserCollection.FindOne(ctx, bson.M{key: value}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return
}

func (s *Store) Insert(ctx context.Context, user User) (*mongo.InsertOneResult, error) {

	insert, err := s.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return insert, nil
}

func (s *Store) Replace(ctx context.Context, user *User, arg interface{}) (*mongo.UpdateResult, error) {
	result, err := s.UserCollection.ReplaceOne(ctx, user, arg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Store) CreateTokens(ctx context.Context, uid string, duration time.Duration) (map[string]string, error) {

	tokens := make(map[string]string)

	access, err := s.Token.CreateAccessToken(uid, duration)
	if err != nil {
		return nil, err
	}

	refresh, err := s.Token.CreateRefreshToken()
	if err != nil {
		return nil, err
	}

	var arg = User{
		UserID:       uid,
		RefreshToken: refresh,
		RefreshExp:   time.Now().Add(time.Hour * 72), //refresh expires after 3 days
	}

	// insert refresh into db
	_, err = s.Insert(ctx, arg)
	if err != nil {
		return nil, err
	}

	tokens["access_token"] = access
	tokens["refresh_token"] = refresh

	return tokens, nil
}

func (user *User) Valid() error {
	if time.Now().After(user.RefreshExp) {
		return errors.New("refresh expired")
	}
	return nil
}
