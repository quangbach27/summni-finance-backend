package errors_test

import (
	"encoding/json"
	common_errors "sumni-finance-backend/internal/common/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Aggregator Tests ---

func TestValidationErrors_Add(t *testing.T) {
	v := &common_errors.ValidationErrors{}

	// 1. Add first error
	v.Add("field1", "error1")
	require.Len(t, v.Errors, 1)
	assert.Equal(t, "field1", v.Errors[0].Field)
	assert.Equal(t, "error1", v.Errors[0].Message)

	// 2. Add second error
	v.Add("field2", "error2")
	require.Len(t, v.Errors, 2)
	assert.Equal(t, "field2", v.Errors[1].Field)
}

func TestValidationErrors_Merge(t *testing.T) {
	t.Run("Merge valid other", func(t *testing.T) {
		v1 := &common_errors.ValidationErrors{}
		v1.Add("f1", "m1")

		v2 := &common_errors.ValidationErrors{}
		v2.Add("f2", "m2")
		v2.Add("f3", "m3")

		v1.Merge(v2)

		require.Len(t, v1.Errors, 3)
		assert.Equal(t, "f1", v1.Errors[0].Field)
		assert.Equal(t, "f2", v1.Errors[1].Field)
		assert.Equal(t, "f3", v1.Errors[2].Field)
	})

	t.Run("Merge nil other", func(t *testing.T) {
		v1 := &common_errors.ValidationErrors{}
		v1.Add("f1", "m1")

		v1.Merge(nil) // Should not panic

		require.Len(t, v1.Errors, 1)
	})

	t.Run("Merge empty other", func(t *testing.T) {
		v1 := &common_errors.ValidationErrors{}
		v1.Add("f1", "m1")

		v2 := &common_errors.ValidationErrors{} // Empty

		v1.Merge(v2)

		require.Len(t, v1.Errors, 1)
	})
}

func TestValidationErrors_AsError(t *testing.T) {
	t.Run("Returns nil when empty", func(t *testing.T) {
		v := &common_errors.ValidationErrors{}
		err := v.AsError()
		assert.Nil(t, err, "AsError() should return nil interface when empty")
	})

	t.Run("Returns error interface when not empty", func(t *testing.T) {
		v := &common_errors.ValidationErrors{}
		v.Add("f1", "m1")

		err := v.AsError()
		assert.NotNil(t, err)

		// Verify type assertion works
		var ve *common_errors.ValidationErrors
		assert.ErrorAs(t, err, &ve)
	})
}

func TestValidationErrors_Error_JSON(t *testing.T) {
	v := &common_errors.ValidationErrors{}
	v.Add("username", "required")
	v.Add("age", "invalid")

	jsonStr := v.Error()

	// 1. Verify it is a valid JSON string
	assert.NotEmpty(t, jsonStr)

	// 2. Unmarshal it back to check structure
	var resultMap map[string][]map[string]string
	err := json.Unmarshal([]byte(jsonStr), &resultMap)
	require.NoError(t, err, "Error() output should be valid JSON")

	// 3. Check contents
	errorsList, exists := resultMap["errors"]
	require.True(t, exists, "JSON should contain 'errors' key")
	require.Len(t, errorsList, 2)

	assert.Equal(t, "username", errorsList[0]["field"])
	assert.Equal(t, "required", errorsList[0]["message"])
	assert.Equal(t, "age", errorsList[1]["field"])
}
