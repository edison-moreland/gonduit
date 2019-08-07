package jwt

// echo.labstack.com/cookbook/jwt

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/edison-moreland/gonduit/models"
)

const jwtUsernameClaim = "name" // Where is the username stored?

// Generate creates and signs a new JWT
func Generate(user *models.User, signingKey []byte, timeToLive time.Duration) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims[jwtUsernameClaim] = user.Username
	//claims["admin"] = false // No admins, yet
	claims["exp"] = time.Now().Add(timeToLive).Unix()

	// Sign token
	tokenSigned, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("Error signing token: %#v", err.Error())
	}

	return tokenSigned, nil
}

// Validate de-signs and parses a JWTString and returns the user associated
func Validate(tokenString string, signingKey []byte) (models.User, error) {
	if IsRevoked(tokenString) {
		return models.User{}, errors.New("Token has been revoked")
	}

	// Decrypt and parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check JWT signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Everything looks good, return signing key
		return signingKey, nil
	})

	// Extract claims, find user
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Grab user name from JWT claims and get user object
		user, err := models.GetUser(claims[jwtUsernameClaim].(string))
		if err != nil {
			return models.User{}, fmt.Errorf("Could not find user %v", claims["user"].(string))
		}

		// JWT valid and user found
		return user, nil
	}
	return models.User{}, fmt.Errorf("Could not validate token. Reason: %v", err.Error())

}
