package validator_test

import (
	"errors"
	"regexp"
	"sumni-finance-backend/internal/common/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: This test file assumes that the Validator and ErrorList structs
// are accessible via the `validator` package and that the ErrorList
// implements the standard Go `error` interface and has the public `Errors` map
// used for assertion checks.

// --- 1. Test Core Functionality and Fluent Interface ---

func TestValidator_Err(t *testing.T) {
	t.Run("Should return error when checks fail", func(t *testing.T) {
		v := validator.New()

		// Test Check and Fluent Interface (fluent methods should return the validator)
		v.Check(false, "field1", "message1").Check(true, "field2", "message2")

		// Test Err() when errors exist
		require.Error(t, v.Err(), "Err() should return an error when checks fail")
	})

	t.Run("Should return nil when no checks fail", func(t *testing.T) {
		v := validator.New()

		// All checks pass
		v.Check(true, "field1", "message1").Check(true, "field2", "message2")

		// Test Err() when no errors exist
		assert.NoError(t, v.Err(), "Err() should return nil when no checks fail")
	})

	t.Run("Should enforce single error per field", func(t *testing.T) {
		v := validator.New()

		v.Check(false, "field1", "first error").
			Check(false, "field1", "second error (should be ignored)")

		// Should only have 1 error entry in the map
		assert.Len(t, v.Errors.Errors, 1, "Expected 1 error for 'field1' due to Add logic")

		// Check that the first message was preserved
		assert.Equal(t, "first error", v.Errors.Errors["field1"], "Expected first error message to be preserved")
	})
}

// --- 2. Test Helper Checkers (Required, MinLength, MaxLength, IsEmail) ---
func TestValidator_Required(t *testing.T) {
	t.Run("Should fail on empty string", func(t *testing.T) {
		v := validator.New()

		v.Required("", "name") // Should fail
		assert.Error(t, v.Err(), "Required failed to detect empty string")
	})

	t.Run("Should pass on non-empty string", func(t *testing.T) {
		v := validator.New()
		v.Required("test", "name") // Should pass
		assert.NoError(t, v.Err(), "Required incorrectly failed for non-empty string")
	})
}

func TestValidator_MinLength(t *testing.T) {
	t.Run("Should fail when value is too short", func(t *testing.T) {
		v := validator.New()
		v.MinLength("abc", "code", 4) // Fail (length 3 < 4)
		assert.Error(t, v.Err(), "MinLength failed to detect short string")
	})

	t.Run("Should pass when value meets minimum length", func(t *testing.T) {
		v := validator.New()
		v.MinLength("abcd", "code", 4) // Pass
		assert.NoError(t, v.Err(), "MinLength incorrectly failed for exact length")
	})
}

func TestValidator_MaxLength(t *testing.T) {
	t.Run("Should fail when value is too long", func(t *testing.T) {
		v := validator.New()
		v.MaxLength("abcd", "code", 3) // Fail (length 4 > 3)
		assert.Error(t, v.Err(), "MaxLength failed to detect long string")
	})

	t.Run("Should pass when value meets maximum length", func(t *testing.T) {
		v := validator.New()
		v.MaxLength("abc", "code", 3) // Pass
		assert.NoError(t, v.Err(), "MaxLength incorrectly failed for exact length")
	})
}

func TestValidator_IsEmail(t *testing.T) {
	tests := []struct {
		email    string
		wantFail bool
	}{
		{
			"test@example.com",
			false,
		},
		{
			"invalid-email",
			true,
		},
		{
			"a@b.c",
			false,
		},
	}

	for _, tt := range tests {
		v := validator.New()
		v.IsEmail(tt.email, "email")

		if tt.wantFail {
			assert.Error(t, v.Err(), "IsEmail should fail for %s", tt.email)
		} else {
			assert.NoError(t, v.Err(), "IsEmail should pass for %s", tt.email)
		}
	}
}

// --- 3. Test TryMerge Functionality ---
func TestValidator_TryMerge(t *testing.T) {
	t.Run("Merge Success: Merges nested errors", func(t *testing.T) {
		// Use require to ensure setup passes
		v := validator.New()
		v.Required("", "base_field")
		require.Error(t, v.Err(), "Base validator setup failed")

		// Create nested error
		nestedV := validator.New()
		nestedV.Required("", "nested_field")
		nestedErr := nestedV.Err()
		require.Error(t, nestedErr, "Nested validator setup failed")

		// Perform the merge
		wasMerged := v.TryMerge(nestedErr)
		assert.True(t, wasMerged, "TryMerge should return true for a validation error")

		// Check if both errors are present (base_field and nested_field)
		assert.Len(t, v.Errors.Errors, 2, "TryMerge failed to merge errors. Expected 2 errors.")

		// Test merging an error that conflicts with an existing one (should preserve original)
		nestedV2 := validator.New()
		nestedV2.Check(false, "base_field", "NEW conflicting message")
		v.TryMerge(nestedV2.Err())

		// The original message should still be there because ErrorList.Add protects it
		assert.Equal(t, "base_field must not be empty", v.Errors.Errors["base_field"], "Merge should preserve original error on conflict")
	})

	t.Run("Merge Failure: Non-validation error", func(t *testing.T) {
		v := validator.New()

		// 1. Create a non-validation error
		fatalErr := errors.New("database connection failed")

		// 2. Perform the merge
		wasMerged := v.TryMerge(fatalErr)
		assert.False(t, wasMerged, "TryMerge should return false for a non-validation error")

		// 3. Check that no errors were added
		assert.NoError(t, v.Err(), "TryMerge incorrectly added errors or merged a non-validation error")
	})

	t.Run("Merge Success: Nil error", func(t *testing.T) {
		v := validator.New()

		// Perform the merge with nil
		wasMerged := v.TryMerge(nil)
		assert.True(t, wasMerged, "TryMerge should return true for a nil error")

		assert.NoError(t, v.Err(), "Validator should not have errors after merging nil")
	})
}

// --- 4. Test Utility Functions ---

func TestValidator_Matches(t *testing.T) {
	t.Run("Should match regex", func(t *testing.T) {
		rx := regexp.MustCompile("^a.*c$")
		assert.True(t, validator.Matches("abc", rx), "Matches failed for a matching string")
	})

	t.Run("Should not match regex", func(t *testing.T) {
		rx := regexp.MustCompile("^a.*c$")
		assert.False(t, validator.Matches("abd", rx), "Matches incorrectly passed for a non-matching string")
	})
}

func TestValidator_In(t *testing.T) {
	list := []string{"red", "green", "blue"}

	t.Run("Should find existing value", func(t *testing.T) {
		assert.True(t, validator.In("green", list...), "In failed for an existing value")
	})

	t.Run("Should not find non-existing value", func(t *testing.T) {
		assert.False(t, validator.In("yellow", list...), "In incorrectly passed for a non-existing value")
	})
}
