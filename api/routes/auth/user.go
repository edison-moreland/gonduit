package auth

import (
	"encoding/json"
	"github.com/edison-moreland/gonduit/api/helpers"
	"github.com/edison-moreland/gonduit/models"
	"net/http"
)

func currentUser(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := r.Context().Value("user").(models.User)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userResponse{User: user}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {

}
