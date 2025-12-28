package validator

import (
	"errors"
	"fmt"
	"regexp"
)

// --- Shared Constants ---
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// --- The Validator Struct ---

// Validator is the checker object. It delegates error storage to ErrorList.
type Validator struct {
	Errors *ErrorList
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{Errors: NewErrorList()}
}

// --- Core API Methods (Fluent & Expressive) ---

// Check adds an error message to the list only if 'ok' is false. Returns for chaining.
func (v *Validator) Check(isValid bool, field, message string) *Validator {
	if !isValid {
		v.Errors.Add(field, message)
	}

	return v
}

// Failed checks if any errors were collected and returns the errorlist as a Go error (or nil).
func (v *Validator) Err() error {
	return v.Errors.AsError()
}

func (v *Validator) TryMerge(err error) bool {
	if err == nil {
		return true // No error to merge
	}

	// Attempt to cast the error to the specific validation type
	var validationErrs *ErrorList

	if errors.As(err, &validationErrs) {
		// Successful cast: merge the errors and return true (merge successful)
		v.Errors.Merge(validationErrs)
		return true
	}

	// Failed cast: the error is not a validation error
	return false
}

// --- Helper Checkers ---

// Required adds an error if the string is empty.
func (v *Validator) Required(value string, field string) *Validator {
	return v.Check(len(value) > 0, field, fmt.Sprintf("%s must not be empty", field))
}

// MinLength adds an error if the string length is less than the minimum.
func (v *Validator) MinLength(value string, field string, n int) *Validator {
	return v.Check(len(value) >= n, field, fmt.Sprintf("%s must be at least %d characters long", field, n))
}

func (v *Validator) MaxLength(value string, field string, n int) *Validator {
	return v.Check(len(value) <= n, field, fmt.Sprintf("%s must be no more than %d characters long", field, n))
}

// IsEmail adds an error if the provided value does not match the EmailRX pattern.
func (v *Validator) IsEmail(value string, field string) *Validator {
	return v.Check(Matches(value, EmailRX), field, "invalid email format")
}

// --- General Utility Functions --- (Keep these simple and static)

// Matches returns true if a string value matches a specific regexp pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// In returns true if a specific value is in a list of strings.
func In(value string, list ...string) bool {
	// ... (logic remains the same)
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}
