package validator

import (
	"encoding/json"
	"fmt"
)

// ErrorList is a collection of validation errors. It implements the Go 'error' interface.
type ErrorList struct {
	// Errors maps the field name (key) to its first error message (value).
	Errors map[string]string `json:"errors"`
}

// New creates a new, empty ErrorList.
func NewErrorList() *ErrorList {
	return &ErrorList{Errors: make(map[string]string)}
}

// Add appends a new error, but only if the field doesn't already have an error.
func (v *ErrorList) Add(field, message string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}

// Merge aggregates errors from another instance into the current list.
// If both lists have an error for the same field, the current list's error is preserved.
func (v *ErrorList) Merge(other *ErrorList) {
	if other == nil || len(other.Errors) == 0 {
		return
	}

	for field, message := range other.Errors {
		// Reuse Add() to respect single-error-per-field logic
		v.Add(field, message)
	}
}

// IsEmpty returns true if the list contains no errors.
func (v *ErrorList) IsEmpty() bool {
	return len(v.Errors) == 0
}

// Error implements the built-in error interface. It marshals the map to a JSON string.
func (v *ErrorList) Error() string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("validation failed with %d errors (json serialization failed)", len(v.Errors))
	}
	return string(bytes)
}

// AsError returns the ErrorList itself as an error if it contains items, otherwise returns nil.
func (v *ErrorList) AsError() error {
	if v.IsEmpty() {
		return nil
	}
	return v
}
