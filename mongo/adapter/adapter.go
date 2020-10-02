package adapter

import (
	"context"
	"fmt"

	"github.com/jpaldi/go-user-api/mongo"
	mongolib "go.mongodb.org/mongo-driver/mongo"
	mongolibopts "go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient establishes a new connecting to mongodb and returns a client wrapper
// which is tailored to usage in the derived-measure/mongo package
func NewClient(uri string) (*ClientAdapter, error) {
	opts := mongolibopts.Client().ApplyURI(uri)
	client, err := mongolib.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("building mongo client: %s", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connecting mongo client: %s", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("pinging mongo server: %s", err)
	}

	return &ClientAdapter{client}, nil
}

// ClientAdapter wraps the mongo lib client to allow easier mocking
type ClientAdapter struct {
	*mongolib.Client
}

// Collection wraps mongolib client functions to provide access to a collection
func (m ClientAdapter) Collection(db, collection string) mongo.Collection {
	return &CollectionAdapter{m.Client.Database(db).Collection(collection)}
}

// CollectionAdapter wraps a mongo lib collection into a Collection.
type CollectionAdapter struct {
	*mongolib.Collection
}

// Insert adds a document to Database
func (c CollectionAdapter) Insert(ctx context.Context, doc interface{}) error {
	_, err := c.Collection.InsertOne(ctx, doc)
	return err
}
