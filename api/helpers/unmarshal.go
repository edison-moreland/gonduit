package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// UnmarshalRequestBody decodes json from a request body into an object
func UnmarshalRequestBody(body io.Reader, object interface{}) error {
	// Read entire body
	rawBody, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("couldn't read body. Reason: %v", err.Error())
	}

	// Unmarshal
	if err = json.Unmarshal(rawBody, &object); err != nil {
		return fmt.Errorf("couldn't unmarshal body. Reason: %v", err.Error())
	}

	return nil
}

// MarshalResponseBody json encodes and object to a response body and sets the status code
func MarshalResponseBody(w http.ResponseWriter, status int, response interface{}) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}
