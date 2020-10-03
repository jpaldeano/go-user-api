package main

import (
	"context"
	"encoding/json"
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

type health struct {
	db *adapter.ClientAdapter
}

func main() {
	ctx := context.Background()
	router := mux.NewRouter()
	database := mustBuildMongoAdapter(ctx)
	mongoDB := mongo.Mongo{
		Client: database.Collection(mongoDatabaseName, mongoCollectionName),
	}

	healthChecker := health{
		db: database,
	}

	mustBuildRoutes(router, mongoDB, healthChecker)

	err := http.ListenAndServe(servicePort, router)
	if err != nil {
		panic(err)
	}
}

func mustBuildRoutes(r *mux.Router, db mongo.Mongo, healthChecker health) {

	usersHandler := handlers.Handler{
		Database: db,
	}
	r.HandleFunc("/health", healthChecker.health).Methods(http.MethodGet)
	r.HandleFunc("/users", usersHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users", usersHandler.GetUsers).Methods(http.MethodGet).Queries()
	r.HandleFunc("/users/{userid}", usersHandler.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{userid}", usersHandler.RemoveUser).Methods(http.MethodDelete)

}

func mustBuildMongoAdapter(ctx context.Context) *adapter.ClientAdapter {
	cl, err := adapter.NewClient(mongoURI)
	if err != nil {
		panic(err)
	}
	return cl
}

func (h health) health(w http.ResponseWriter, r *http.Request) {
	var databaseStatus = "OK"
	err := h.db.Ping(context.TODO(), nil)
	if err != nil {
		databaseStatus = "UNHEALTHY"
	}

	type HealthStatusResponse struct {
		Database string `json:"database_status"`
	}

	response := HealthStatusResponse{
		Database: databaseStatus,
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}
