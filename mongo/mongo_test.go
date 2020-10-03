package mongo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jpaldi/go-user-api/mongo"
	mongolib "go.mongodb.org/mongo-driver/mongo"
)

type mockClient struct {
	client mockDatabase
}
type mockDatabase struct {
	insertOne        func(ctx context.Context, doc interface{}) error
	findOneAndUpdate func(ctx context.Context, filter interface{}, update interface{}) *mongolib.SingleResult
	deleteOne        func(ctx context.Context, filter interface{}) error
	find             func(ctx context.Context, query interface{}) (*mongolib.Cursor, error)
}

func (m mockDatabase) InsertOne(ctx context.Context, doc interface{}) error {
	return m.insertOne(ctx, doc)
}

func (m mockDatabase) DeleteOne(ctx context.Context, filter interface{}) error {
	return m.deleteOne(ctx, filter)
}

func (m mockDatabase) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) *mongolib.SingleResult {
	return m.findOneAndUpdate(ctx, filter, update)
}

func (m mockDatabase) Find(ctx context.Context, query interface{}) (*mongolib.Cursor, error) {
	return m.find(ctx, query)
}

func mockInsertDatabaseFailure() mockDatabase {
	return mockDatabase{
		insertOne: func(ctx context.Context, doc interface{}) error {
			return fmt.Errorf("database error on creating a user")
		},
	}
}

func mockInsertDatabaseOK() mockDatabase {
	return mockDatabase{
		insertOne: func(ctx context.Context, doc interface{}) error {
			return nil
		},
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name                  string
		database              mockDatabase
		expectedErrorResponse error
		expectedUser          *mongo.User
	}{
		{
			name:                  "if database adapter returns an error, it should return an error",
			database:              mockInsertDatabaseFailure(),
			expectedErrorResponse: fmt.Errorf("cannot insert: database error on creating a user"),
		},
		{
			name:     "if database adapter inserts successfuly, it should return the new user",
			database: mockInsertDatabaseOK(),
			expectedUser: &mongo.User{
				Nickname: "test",
			},
		},
	} {

		t.Run(tt.name, func(t *testing.T) {
			client := mongo.Mongo{
				Client: tt.database,
			}

			nickname := "test"
			usr, err := client.CreateUser(context.Background(), nickname, "", "", "", "", "")

			if usr != nil {
				if usr.Nickname != tt.expectedUser.Nickname {
					// I am just testing Nickname because I don't want to go much in detail and spend much time here.
					// If I want to compare objects I would use this library: github.com/google/go-cmp/cmp
					t.Fatalf("wrong Nickname: got %s want %s",
						usr.Nickname, tt.expectedUser.Nickname)
				}
			}

			if err != nil {
				if err.Error() != tt.expectedErrorResponse.Error() {
					t.Fatalf("wrong error: got %s want %s",
						err.Error(), tt.expectedErrorResponse.Error())
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// TODO
}

func TestRemoveUser(t *testing.T) {
	// TODO
}

func TestGetUsers(t *testing.T) {
	// TODO
}
