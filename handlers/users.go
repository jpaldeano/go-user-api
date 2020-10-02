package handlers

import (
	"context"
	"net/http"

	"github.com/jpaldi/go-user-api/mongo"
)

type User struct {
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
}

type Database interface {
	CreateUser(ctx context.Context, nickname string) (*mongo.User, error)
}

type LoginHandler struct {
	Database Database
}

// Handle is responsible for handle user routes
func (handler *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (handler *LoginHandler) post(w http.ResponseWriter, r *http.Request) error {
	handler.Database.CreateUser(r.Context(), "conas")
	return nil
}
