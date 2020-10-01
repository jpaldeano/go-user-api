package handlers

import (
	"encoding/json"
	"net/http"
)

func User struct {
	FirstName string
	LastName string
	Nickname string
	Password string
	Email string
	Country string
}

type LoginHandler struct {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)

	return nil
}
