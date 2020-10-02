package mongo

import (
	"context"
	"fmt"
)

// User represents the object stored in database.
type User struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Nickname  string `json:"nickname" bson:"nickname"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Password  string `json:"password" bson:"password"`
	Email     string `json:"email" bson:"email"`
	Country   string `json:"country" bson:"country"`
}

// Collection represents the interface to wrap the collection from mongo drive
type Collection interface {
	InsertOne(ctx context.Context, doc interface{}) error
}

// Mongo represents a mongo client wrapped to provide service-specific functionality.
type Mongo struct {
	Client Collection
}

// CreateUser creates a user and returns the object if is successfully inserted
func (mgo Mongo) CreateUser(ctx context.Context, nickname string, firstname string, lastname string, password string, email string, country string) (*User, error) {
	user := User{
		Nickname:  nickname,
		FirstName: firstname,
		LastName:  lastname,
		Password:  password,
		Email:     email,
		Country:   country,
	}

	if err := mgo.Client.InsertOne(ctx, user); err != nil {
		return nil, fmt.Errorf("cannot insert: %s", err)
	}

	return &user, nil
}
