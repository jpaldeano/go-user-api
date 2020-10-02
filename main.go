package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	router := mux.NewRouter()

	mustBuildRoutes(router, database)

	err := http.ListenAndServe(servicePort, router)
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
}

func mustBuildRoutes(r *mux.Router, db *adapter.ClientAdapter) {
	usersHandler := handlers.Handler{
		Database: mongo.Mongo{
			Client: db.Collection(mongoDatabaseName, mongoCollectionName),
		},
	}
	r.HandleFunc("/users", usersHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{userid}", usersHandler.UpdateUser).Methods(http.MethodPut)
}

func mustBuildMongoAdapter(ctx context.Context) *adapter.ClientAdapter {
	cl, err := adapter.NewClient(mongoURI)
	if err != nil {
		panic(err)
	}

	return cl
}
