package authentication

// echo.labstack.com/cookbook/jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/edison-moreland/gonduit/models"
)

func GenerateJWT(user models.User, signingKey []byte, lifetime time.Duration) (jwt.Token, error) {
	// Create token
	token := jwt.New(jwt.Sign)
}
