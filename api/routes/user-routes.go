package routes

import (
	"github.com/edison-moreland/gonduit/api/helpers"
	"github.com/edison-moreland/gonduit/authentication"
	"github.com/edison-moreland/gonduit/authentication/jwt"
	"github.com/edison-moreland/gonduit/models"
	"github.com/gorilla/mux"
	"net/http"
)

func AddUserRoutes(router *mux.Router) {
	router.Path("/users/login").Methods(http.MethodPost).HandlerFunc(login).Name("login")
	router.Path("/users").Methods(http.MethodPost).HandlerFunc(register).Name("register")
	router.Path("/user").Methods(http.MethodGet).Handler(jwt.Required(currentUser)).Name("currentuser")
	router.Path("/user").Methods(http.MethodPut).Handler(jwt.Required(updateUser)).Name("updateuser")
}

type userResponse struct {
	User models.User `json:"user"`
}

func login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Expected body
	loginBody := struct {
		User struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		} `json:"user" validate:"required"`
	}{}

	err := helpers.UnmarshalRequestBody(r.Body, &loginBody)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	err = helpers.ValidateStruct(loginBody)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// Authenticate login, returns user model + JWT
	user, err := authentication.Login(loginBody.User.Username, loginBody.User.Password)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// Return user object
	if err := helpers.MarshalResponseBody(w, http.StatusOK, userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Expected body
	registerBody := struct {
		User struct {
			Username string `json:"username" validate:"required"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
		} `json:"user" validate:"required"`
	}{}

	err := helpers.UnmarshalRequestBody(r.Body, &registerBody)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	err = helpers.ValidateStruct(registerBody)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// Create new user
	newUser := models.User{Username: registerBody.User.Username, Email: registerBody.User.Email}
	err = newUser.UpdatePassword(registerBody.User.Password)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
	newUser.Save()

	// Login to get JWT
	user, err := authentication.Login(registerBody.User.Username, registerBody.User.Password)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// Return new user
	if err := helpers.MarshalResponseBody(w, http.StatusOK, userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func currentUser(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := jwt.CurrentUser(r)

	if err := helpers.MarshalResponseBody(w, http.StatusOK, userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Get current user from context
	user := jwt.CurrentUser(r)

	// Expected body
	updateBody := struct {
		User struct {
			Username string `json:"username"`
			Email    string `json:"email" validate:"omitempty,email"`
			Bio      string `json:"bio"`
			Image    string `json:"image" validate:"omitempty,url"`
			Password string `json:"password"`
		} `json:"user" validate:"required"`
	}{}

	err := helpers.UnmarshalRequestBody(r.Body, &updateBody)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	err = helpers.ValidateStruct(updateBody)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// If password field given, update password
	if updateBody.User.Password != "" {
		err := user.UpdatePassword(updateBody.User.Password)
		if err != nil {
			helpers.Err422(err.Error(), w)
			return
		}
	}

	// Shove every other field into a user object and merge it to the original
	user.UpdateUser(models.User{
		Username: updateBody.User.Username,
		Email:    updateBody.User.Email,
		Bio:      updateBody.User.Bio,
		Image:    updateBody.User.Image,
	})

	// Save any changes made
	user.Save()

	// Return updated user
	w.WriteHeader(http.StatusOK)
	if err := helpers.MarshalResponseBody(w, http.StatusOK, userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}
