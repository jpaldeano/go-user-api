package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/jpaldi/go-user-api/handlers"
)

func main() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("api running")
}
func mustStartRoutes() {
	loginHandler := handlers.LoginHandler{}
	http.HandleFunc("/users", loginHandler.Handle)

}

func mustLoadConfig() {
	godotenv.Load("secrets.env")
}

func mustBuildMongoAdapter() {

}
