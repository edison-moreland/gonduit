package auth

import (
	"encoding/json"
	"github.com/edison-moreland/gonduit/api/helpers"
	"github.com/edison-moreland/gonduit/authentication"
	"github.com/edison-moreland/gonduit/authentication/jwt"
	"github.com/edison-moreland/gonduit/models"
	"github.com/gorilla/mux"
	"net/http"
)

func AddRoutes(router *mux.Router) {
	router.Path("/users/login").Methods(http.MethodPost).HandlerFunc(login).Name("login")
	router.Path("/users").Methods(http.MethodPost).HandlerFunc(register).Name("register")
	router.Path("/user").Methods(http.MethodGet).Handler(jwt.JWTRequired(currentUser)).Name("currentuser")
	router.Path("/user").Methods(http.MethodPut).Handler(jwt.JWTRequired(updateUser)).Name("updateuser")
}

func login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request := userRequest{User: userRequestBody{}}
	err := helpers.UnmarshalRequestBody(r.Body, &request)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	user, err := authentication.Login(request.User.Username, request.User.Password)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request := userRequest{User: userRequestBody{}}
	err := helpers.UnmarshalRequestBody(r.Body, &request)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// Make sure email field is filled (its optional on struct)
	if request.User.Email == "" {
		helpers.Err422("Email field can't be blank", w)
		return
	}

	// Create user
	newUser := models.User{Username: request.User.Username, Email: request.User.Email}
	err = newUser.UpdatePassword(request.User.Password)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
	newUser.Save()

	// Login
	user, err := authentication.Login(request.User.Username, request.User.Password)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func currentUser(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := r.Context().Value("user").(models.User)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := r.Context().Value("user").(models.User)
}
