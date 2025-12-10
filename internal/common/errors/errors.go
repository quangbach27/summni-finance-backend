package errors

import (
	"encoding/json"
	"fmt"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// --- The Aggregator ---

type ValidationErrors struct {
	Errors []*ValidationError `json:"errors"`
}

// Add appends a new error
func (v *ValidationErrors) Add(field, message string) {
	v.Errors = append(v.Errors, &ValidationError{
		Field:   field,
		Message: message,
	})
}

// Merge aggregates errors from another instance
func (v *ValidationErrors) Merge(other *ValidationErrors) {
	if other == nil || len(other.Errors) == 0 {
		return
	}
	v.Errors = append(v.Errors, other.Errors...)
}

// AsError (formerly Return) returns nil if empty
func (v *ValidationErrors) AsError() error {
	if len(v.Errors) == 0 {
		return nil
	}
	return v
}

// Error implements the standard error interface.
// Now it returns a JSON string representing the object.
func (v *ValidationErrors) Error() string {
	// Marshal the entire struct (which contains the "errors" list)
	bytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("validation failed with %d errors (json serialization failed)", len(v.Errors))
	}
	return string(bytes)
}
