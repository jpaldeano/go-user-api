package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jpaldi/go-user-api/mongo"
)

// UsersDatabase wraps the Database client functions
type UsersDatabase interface {
	CreateUser(ctx context.Context, nickname string) (*mongo.User, error)
}

// UsersHandler represents the handler for users routes
type UsersHandler struct {
	Database UsersDatabase
}

// HandleUsers is responsible for handle user routes
func (handler *UsersHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
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

func (handler *UsersHandler) post(w http.ResponseWriter, r *http.Request) error {
	userBody := &userRequestBody{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(userBody); err != nil {
		panic(err)
	}

	if validErrs := userBody.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		writeResponse(w, http.StatusBadRequest, err)
	}

	// todo: Hash & Salt passwords: https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
	user, err := handler.Database.CreateUser(r.Context(), userBody.Nickname)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, err)
	}

	// In case User, was inserted return the user object
	writeResponse(w, http.StatusOK, user)
	return nil
}

func writeResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
