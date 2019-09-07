package jwt

import (
	"context"
	"log"
	"net/http"

	"github.com/edison-moreland/gonduit/models"
)

// Create new type for use as context key to avoid collisions
// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
type key int

const userKey key = iota

// Required ensures token in request and uses token to get current user
func Required(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token
		user, err := ValidateFromRequest(r)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Add current user to context
		ctx := context.WithValue(r.Context(), userKey, user)

		next(w, r.WithContext(ctx))
	})
}

// Optional adds the current user to the request context if logged in, otherwise adds an empty user
func Optional(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to validate token
		user, err := ValidateFromRequest(r)
		if err != nil {
			// No user logged in, add empty user to context
			ctx := context.WithValue(r.Context(), userKey, models.User{})

			next(w, r.WithContext(ctx))

			return
		}

		// Add current user to context
		ctx := context.WithValue(r.Context(), userKey, user)

		next(w, r.WithContext(ctx))
	})
}

// CurrentUser returns the logged in user from the request context
func CurrentUser(r *http.Request) models.User {
	return r.Context().Value(userKey).(models.User)
}

// UserLoggedIn returns true if a valid user model exists in the request context
func UserLoggedIn(r *http.Request) bool {
	return CurrentUser(r).Username != ""
}
