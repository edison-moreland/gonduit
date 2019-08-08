package jwt

import (
	"context"
	"log"
	"net/http"
)

// JWTRequired ensures token in request and uses token to get current user
func JWTRequired(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token
		user, err := ValidateFromRequest(r)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Add current user to context
		ctx := context.WithValue(r.Context(), "user", user)

		next(w, r.WithContext(ctx))
	})
}
