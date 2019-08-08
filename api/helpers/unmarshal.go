package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

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
