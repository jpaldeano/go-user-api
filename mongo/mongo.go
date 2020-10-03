package mongo

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongolib "go.mongodb.org/mongo-driver/mongo"
)

var (
	validURLParams = []string{"nickname", "first_name", "country", "last_name", "email"}
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

// Collection represents the interface to wrap the mongo drive collection
type Collection interface {
	InsertOne(ctx context.Context, doc interface{}) error
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) *mongolib.SingleResult
	DeleteOne(ctx context.Context, filter interface{}) error
	Find(ctx context.Context, query interface{}) (*mongolib.Cursor, error)
}

// Mongo represents a mongo client wrapped to provide service-specific functionality.
type Mongo struct {
	Client Collection
}

// CreateUser creates a user and returns the object if is successfully inserted
func (mgo Mongo) CreateUser(ctx context.Context, nickname string, firstname string, lastname string, password string, email string, country string) (*User, error) {
	user := User{
		ID:        uuid.New().String(),
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

// UpdateUser creates a user and returns the object if is successfully inserted
func (mgo Mongo) UpdateUser(ctx context.Context, guid string, nickname string, firstname string, lastname string, password string, email string, country string) (*User, error) {

	update := bson.M{
		"$set": bson.M{
			"nickname":   nickname,
			"first_name": firstname,
			"last_name":  lastname,
			"password":   password,
			"email":      email,
			"country":    country,
		},
	}
	result := mgo.Client.FindOneAndUpdate(ctx, bson.M{"_id": guid}, update)
	if result.Err() != nil {
		return nil, result.Err()
	}

	doc := bson.M{}
	decodeErr := result.Decode(&doc)
	if decodeErr != nil {
		return nil, decodeErr
	}

	user := User{}
	bsonBytes, _ := bson.Marshal(doc)

	err := bson.Unmarshal(bsonBytes, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// RemoveUser removes a user from mongo
func (mgo Mongo) RemoveUser(ctx context.Context, guid string) error {
	filter := bson.M{
		"_id": guid,
	}
	err := mgo.Client.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// GetUsers get a users from mongo
func (mgo Mongo) GetUsers(ctx context.Context, params url.Values) ([]*User, error) {
	query := bson.M{}
	for k, v := range params {
		// check if the query parameter is expected to avoid SQL Injections
		// for the purpose of this service, assume this route only allows one single parameter per key
		if contains(validURLParams, k) {
			query[k] = v[0]
		}
	}

	cursor, err := mgo.Client.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	users := []*User{}

	for cursor.Next(ctx) {
		doc := bson.M{}
		if err = cursor.Decode(&doc); err != nil {
			return nil, err
		}

		bytes, err := bson.Marshal(doc)
		if err != nil {
			return nil, err
		}

		u := User{}
		err = bson.Unmarshal(bytes, &u)
		if err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	return users, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
