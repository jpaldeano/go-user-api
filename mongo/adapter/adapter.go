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

// InsertOne adds a document in Mongo
func (c CollectionAdapter) InsertOne(ctx context.Context, doc interface{}) error {
	_, err := c.Collection.InsertOne(ctx, doc)
	return err
}

// FindOneAndUpdate and updates a document to Database
func (c CollectionAdapter) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) *mongolib.SingleResult {
	after := mongolibopts.After
	opt := mongolibopts.FindOneAndUpdateOptions{
		ReturnDocument: &after, // ReturnDocument option to return the updated document
	}
	return c.Collection.FindOneAndUpdate(ctx, filter, update, &opt)
}

// DeleteOne removes a document from Mongo
func (c CollectionAdapter) DeleteOne(ctx context.Context, filter interface{}) (*mongolib.DeleteResult, error) {
	return c.Collection.DeleteOne(ctx, filter)
}

// Find returns all documents from Mongo matching the given query
func (c CollectionAdapter) Find(ctx context.Context, query interface{}) (*mongolib.Cursor, error) {
	return c.Collection.Find(ctx, query)
}
