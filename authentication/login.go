package authentication

import (
	"errors"
	"fmt"

	"github.com/edison-moreland/gonduit/authentication/jwt"
	"github.com/edison-moreland/gonduit/models"
)

// Login compares username and password, return user object with token if they match
func Login(username, password string) (models.User, error) {
	// Find user in database
	user, err := models.GetUser(username)
	if err != nil {
		// TODO: Upstream error is getting lost, log it
		return models.User{}, errors.New("unknown username/password")
	}

	// Check password
	if !user.CheckPassword(password) {
		// Password didnt match
		return models.User{}, errors.New("unknown username/password")
	}

	// User logged in, generate token
	token, err := jwt.Generate(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("error generating token: %v", token)
	}

	user.Token = token
	return user, nil
}

// Logout revokes a users JWT
func Logout(token string) {
	// This might have more logic later
	jwt.Revoke(token)
}
