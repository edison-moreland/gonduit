package helpers

import (
	"gopkg.in/go-playground/validator.v9"
)

var _validator *validator.Validate
var _validatorIsInitialized = false

// ValidateStruct uses validator library to validate a struct
func ValidateStruct(unvalidated interface{}) error {
	if !_validatorIsInitialized {
		// Only create validator once, who knows if this is thread safe
		// Opens the door for custom validators later down the line
		_validator = validator.New()
		_validatorIsInitialized = true
	}

	// ValidateStruct struct
	return _validator.Struct(unvalidated)
}
