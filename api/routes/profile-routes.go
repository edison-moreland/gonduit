package routes

import (
	"github.com/edison-moreland/tokenware"
	"net/http"

	"github.com/edison-moreland/gonduit/api/helpers"
	"github.com/edison-moreland/gonduit/models"
	"github.com/gorilla/mux"
)

// AddProfileRoutes adds all user related routes to a gorilla router
func AddProfileRoutes(router *mux.Router) {
	router.Path("/profile/{username}").Methods(http.MethodGet).Handler(tokenware.Optional(getProfile)).Name("getprofile")
	router.Path("/profile/{username}/follow").Methods(http.MethodPost).Handler(tokenware.Required(followUser)).Name("followUser")
	router.Path("/profile/{username}/follow").Methods(http.MethodDelete).Handler(tokenware.Required(unfollowUser)).Name("unfollowUser")
}

type profileResponse struct {
	Profile models.Profile `json:"profile"`
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	// Get username from url
	username := mux.Vars(r)["username"]

	// Find user
	profileUser, err := models.GetUser(username)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// If user is logged in, find out if they're following this profile
	following := false
	current := helpers.CurrentUser(r)
	if current.Username != "" {
		following = current.IsFollowingUser(username)
	}

	// Generate profile
	profile := profileUser.GetProfile(following)

	// Write it out
	if err = helpers.MarshalResponseBody(w, http.StatusOK, profileResponse{Profile: profile}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func followUser(w http.ResponseWriter, r *http.Request) {
	// Get username from url
	username := mux.Vars(r)["username"]

	// Find user
	profileUser, err := models.GetUser(username)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// Follow user
	loggedInUser := helpers.CurrentUser(r)
	err = loggedInUser.FollowUser(username)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
	loggedInUser.Save()

	// Generate profile
	profile := profileUser.GetProfile(loggedInUser.IsFollowingUser(username))

	// Write it out
	if err = helpers.MarshalResponseBody(w, http.StatusOK, profileResponse{Profile: profile}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}

func unfollowUser(w http.ResponseWriter, r *http.Request) {
	// Find user with username from url
	profileUser, err := models.GetUser(mux.Vars(r)["username"])
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}

	// Follow user
	loggedInUser := helpers.CurrentUser(r)
	err = loggedInUser.UnFollowUser(profileUser.Username)
	if err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
	loggedInUser.Save()

	// Generate profile
	profile := profileUser.GetProfile(loggedInUser.IsFollowingUser(profileUser.Username))

	// Write it out
	if err = helpers.MarshalResponseBody(w, http.StatusOK, profileResponse{Profile: profile}); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}
