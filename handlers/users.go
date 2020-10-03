package handlers

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/jpaldi/go-user-api/mongo"
)

// UsersDatabase wraps the Database client functions
type UsersDatabase interface {
	CreateUser(ctx context.Context, nickname string, firstname string, lastname string, password string, email string, country string) (*mongo.User, error)
	UpdateUser(ctx context.Context, guid string, nickname string, firstname string, lastname string, password string, email string, country string) (*mongo.User, error)
	RemoveUser(ctx context.Context, guid string) error
	GetUsers(ctx context.Context, params url.Values) ([]*mongo.User, error)
}

// Handler represents the handler for users routes
type Handler struct {
	Database UsersDatabase
}

// CreateUser handles the POST /users request
func (handler *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := validateJSON(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if validErrs := userBody.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		writeResponse(w, http.StatusBadRequest, err)
		return
	}

	user, err := handler.Database.CreateUser(r.Context(), userBody.Nickname, userBody.FirstName, userBody.LastName, userBody.Password, userBody.Email, userBody.Country)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	// In case User, was inserted return the user object
	writeResponse(w, http.StatusOK, user)
}

// UpdateUser handles the Put /users/{userid} request
func (handler *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userid := mux.Vars(r)["userid"]

	userBody, err := validateJSON(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if validErrs := userBody.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		writeResponse(w, http.StatusBadRequest, err)
		return
	}

	user, err := handler.Database.UpdateUser(r.Context(), userid, userBody.Nickname, userBody.FirstName, userBody.LastName, userBody.Password, userBody.Email, userBody.Country)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err)
		return
	}

	// In case User, was inserted return the user object
	writeResponse(w, http.StatusOK, user)
}

// RemoveUser handles the DELETE /users/{userid} request
func (handler *Handler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	userid := mux.Vars(r)["userid"]

	err := handler.Database.RemoveUser(r.Context(), userid)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err)
		return
	}

	// In case User, was inserted return the user object
	writeResponse(w, http.StatusOK, nil)
}

// GetUsers handles the GET /users request
func (handler *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	results, err := handler.Database.GetUsers(r.Context(), queryParams)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	// In case User, was inserted return the user object
	writeResponse(w, http.StatusOK, results)
}
