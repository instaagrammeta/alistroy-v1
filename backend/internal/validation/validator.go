package validation

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

// Struct validates a struct using go-playground/validator and returns an error if any.
func Struct(s any) error {
	return v.Struct(s)
}

// Validator returns the underlying validator instance for advanced registration.
func Validator() *validator.Validate { return v }
