package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jpaldi/go-user-api/handlers"
	"github.com/jpaldi/go-user-api/mongo"
	adapter "github.com/jpaldi/go-user-api/mongo/adapter"
)

var (
	mongoDatabaseName   = os.Getenv("MONGO_DATABASE_NAME")
	mongoCollectionName = os.Getenv("MONGO_COLLECTION_NAME")
)

func main() {
	mustLoadConfig()

	ctx := context.Background()
	database := mustBuildMongoAdapter(ctx)
	mustBuildRoutes(database)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("api running")
}

func mustBuildRoutes(db *adapter.ClientAdapter) {
	loginHandler := handlers.LoginHandler{
		Database: mongo.Mongo{
			Client: db.Collection(mongoDatabaseName, mongoCollectionName),
		},
	}
	http.HandleFunc("/users", loginHandler.Handle)
}

func mustLoadConfig() {
	godotenv.Load("secrets.env")
}

func mustBuildMongoAdapter(ctx context.Context) *adapter.ClientAdapter {
	cl, err := adapter.NewClient(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}

	return cl
}
