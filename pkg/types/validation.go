package types

import (
	"fmt"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	
	msg := fmt.Sprintf("%d validation error(s):", len(ve))
	for i, err := range ve {
		msg += fmt.Sprintf("\n%d. %s", i+1, err.Error())
	}
	return msg
}

// Add adds a validation error to the collection
func (ve *ValidationErrors) Add(field, message string, value interface{}) {
	*ve = append(*ve, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors returns true if there are validation errors
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}