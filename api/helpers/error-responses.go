package helpers

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Errors apiErrorBody `json:"errors"`
}

type apiErrorBody struct {
	Body []string `json:"body"`
}

// Err422 writes an error in the standard format to the response body
func Err422(reason string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)

	errResponse := apiError{
		Errors: apiErrorBody{Body: []string{reason}},
	}

	_ = json.NewEncoder(w).Encode(errResponse)
}
