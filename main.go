package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jpaldi/go-user-api/handlers"
	"github.com/jpaldi/go-user-api/mongo"
	adapter "github.com/jpaldi/go-user-api/mongo/adapter"
)

var (
	servicePort         = fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT"))
	mongoURI            = os.Getenv("MONGO_URI")
	mongoDatabaseName   = os.Getenv("MONGO_DATABASE_NAME")
	mongoCollectionName = os.Getenv("MONGO_COLLECTION_NAME")
)

func main() {
	ctx := context.Background()
	database := mustBuildMongoAdapter(ctx)

	mustBuildRoutes(database)

	err := http.ListenAndServe(servicePort, nil)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}

func mustBuildRoutes(db *adapter.ClientAdapter) {
	usersHandler := handlers.UsersHandler{
		Database: mongo.Mongo{
			Client: db.Collection(mongoDatabaseName, mongoCollectionName),
		},
	}
	http.HandleFunc("/users", usersHandler.HandleUsers)
}

func mustBuildMongoAdapter(ctx context.Context) *adapter.ClientAdapter {
	fmt.Println(mongoURI)
	cl, err := adapter.NewClient(mongoURI)
	if err != nil {
		panic(err)
	}

	return cl
}
