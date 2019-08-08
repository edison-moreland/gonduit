package auth

import (
	"github.com/edison-moreland/gonduit/authentication/jwt"
	"github.com/gorilla/mux"
	"net/http"
)

func AddRoutes(router *mux.Router) {
	router.Path("/users/login").Methods(http.MethodPost).HandlerFunc(login).Name("login")
	router.Path("/users").Methods(http.MethodPost).HandlerFunc(register).Name("register")
	router.Path("/user").Methods(http.MethodGet).Handler(jwt.JWTRequired(currentUser)).Name("currentuser")
	router.Path("/user").Methods(http.MethodPut).Handler(jwt.JWTRequired(updateUser)).Name("updateuser")
}
