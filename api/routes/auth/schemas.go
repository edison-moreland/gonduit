package auth

import "github.com/edison-moreland/gonduit/models"

type userResponse struct {
	User models.User `json:"user"`
}

type userRequest struct {
	User userRequestBody `json:"user"`
}

type userRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username, omitempty"`
}
