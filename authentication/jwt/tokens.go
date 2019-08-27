package jwt

// echo.labstack.com/cookbook/jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/edison-moreland/gonduit/models"
)

const jwtUsernameClaim = "name" // Where is the username stored?
const jwtSigningKey = "SupoerSecret"
const jwtTimeToLive = time.Hour * 24
const jwtHeader = "Authorization"
const jwtPrefix = "Token "

// Generate creates and signs a new JWT
func Generate(user *models.User) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims[jwtUsernameClaim] = user.Username
	//claims["admin"] = false // No admins, yet
	claims["exp"] = time.Now().Add(jwtTimeToLive).Unix()

	// Sign token
	tokenSigned, err := token.SignedString([]byte(jwtSigningKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %#v", err.Error())
	}

	return tokenSigned, nil
}

// Validate de-signs and parses a JWTString and returns the user associated
func Validate(tokenString string) (models.User, error) {
	if IsRevoked(tokenString) {
		return models.User{}, errors.New("token has been revoked")
	}

	// Decrypt and parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check JWT signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Everything looks good, return signing key
		return []byte(jwtSigningKey), nil
	})

	if err != nil {
		return models.User{}, fmt.Errorf("could not validate token (%v). Reason: %v", tokenString, err.Error())
	}

	// Extract claims, find user
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Grab user name from JWT claims and get user object
		user, err := models.GetUser(claims[jwtUsernameClaim].(string))
		if err != nil {
			return models.User{}, fmt.Errorf("could not find user %v", claims["user"].(string))
		}

		// Add token to user
		user.Token = tokenString

		// JWT valid and user found
		return user, nil
	}
	return models.User{}, fmt.Errorf("could not validate token (%v)", tokenString)

}

// ValidateFromRequest gets token from request headers and validates it
func ValidateFromRequest(r *http.Request) (models.User, error) {
	rawToken := r.Header.Get(jwtHeader)
	if !strings.HasPrefix(rawToken, jwtPrefix) {
		return models.User{}, errors.New("token not in headers")
	}

	token := strings.TrimPrefix(rawToken, jwtPrefix)

	return Validate(token)
}
