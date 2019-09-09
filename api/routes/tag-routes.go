package routes

import (
	"net/http"

	"github.com/edison-moreland/gonduit/api/helpers"
	"github.com/edison-moreland/gonduit/models"
	"github.com/gorilla/mux"
)

// AddTagRoutes adds all tag related routes to a gorilla router
func AddTagRoutes(router *mux.Router) {
	router.Path("/tags").Methods(http.MethodGet).HandlerFunc(getTags).Name("tag")
}

func getTags(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Tags []models.Tag `json:"tags"`
	}{
		Tags: models.GetTags(),
	}

	if err := helpers.MarshalResponseBody(w, 200, response); err != nil {
		helpers.Err422(err.Error(), w)
		return
	}
}
