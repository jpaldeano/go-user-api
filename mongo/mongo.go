package mongo

import (
	"context"
	"fmt"
)

// User represents the object stored in database.
type User struct {
	Nickname string `json:"nickname" bson:"nickname"`
}

// Collection represents the interface to wrap the collection from mongo drive
type Collection interface {
	Insert(ctx context.Context, doc interface{}) error
}

// Database represents a mongo client wrapped to provide service-specific functionality.
type Mongo struct {
	Client Collection
}

// CreateUser creates a user and returns the object if is successfully inserted
func (mgo Mongo) CreateUser(ctx context.Context, nickname string) (*User, error) {
	user := User{
		Nickname: nickname,
	}

	if err := mgo.Client.Insert(ctx, user); err != nil {
		return nil, fmt.Errorf("cannot insert: %s", err)
	}

	return &user, nil
}
