package helpers

import (
	"github.com/edison-moreland/gonduit/models"
	"github.com/edison-moreland/tokenware"
	"net/http"
)

// CurrentUser returns the user associated with the identity in the request context
func CurrentUser(r *http.Request) models.User {
	identity := tokenware.CurrentIdentity(r)

	if identity != nil {
		// Identity was found! get user model and add token
		user, _ := models.GetUser(identity.(string))

		// Ignore error, if this function fails there is no way we got here
		token, _ := tokenware.GetRawToken(r)
		user.Token = token
		return user
	}

	// No identity found, return empty user
	// We don't return an error because in this situation it isn't an error
	// in the context of tokenware.Required() the identity is always added
	// in the context of tokenware.Optional() the identity is expected to sometimes be blank

	return models.User{}
}
