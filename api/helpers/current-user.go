package helpers

import (
	"github.com/edison-moreland/gonduit/models"
	"github.com/edison-moreland/tokenware"
	"net/http"
)

func currentUser(r *http.Request) models.User {
	identity := tokenware.CurrentIdentity(r)

	if identity != nil {
		// Identity was found! get user model
		user, _ = models.GetUser(identity.(string))
		return user
	}

	// No identity found, return empty user
	// We don't return an error because in this situation it isn't an error
	// in the context of tokenware.Required() we al;ready check that an identity is added

	return models.User{}
}
