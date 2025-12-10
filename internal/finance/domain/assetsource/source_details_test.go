package assetsource_test

import (
	"errors"
	"testing"

	commons_errors "sumni-finance-backend/internal/common/errors"
	"sumni-finance-backend/internal/finance/domain/assetsource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBankDetails(t *testing.T) {
	tests := []struct {
		name               string
		inputName          string
		inputAccount       string
		expectError        bool
		expectedErrorCount int
		expectedFields     []string // Fields expected to have errors
	}{
		{
			name:         "Success: Valid Inputs",
			inputName:    "TechBank",
			inputAccount: "987654",
			expectError:  false,
		},
		{
			name:               "Failure: Missing Bank Name Only",
			inputName:          "",
			inputAccount:       "123",
			expectError:        true,
			expectedErrorCount: 1,
		},
		{
			name:               "Failure: Missing Account Number Only",
			inputName:          "FinBank",
			inputAccount:       "",
			expectError:        true,
			expectedErrorCount: 1,
		},
		{
			name:               "Failure: Both Missing (Aggregated)",
			inputName:          "",
			inputAccount:       "",
			expectError:        true,
			expectedErrorCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details, err := assetsource.NewBankDetails(tt.inputName, tt.inputAccount)

			if tt.expectError {
				require.Error(t, err)
				assert.True(t, details.IsZero(), "Should return zero value on error")

				// --- Critical Error Aggregation Checks ---
				var valErrs *commons_errors.ValidationErrors
				require.True(t, errors.As(err, &valErrs), "Error must be of type *ValidationErrors")

				assert.Len(t, valErrs.Errors, tt.expectedErrorCount, "Incorrect number of aggregated errors")
			} else {
				require.NoError(t, err)
				assert.False(t, details.IsZero(), "Should not be zero value on success")
				assert.Equal(t, tt.inputName, details.BankName())
				assert.Equal(t, tt.inputAccount, details.AccountNumber())
			}
		})
	}
}

func TestBankDetails_IsZero(t *testing.T) {
	// Case 1: Zero value (uninitialized)
	var zeroDetails assetsource.BankDetails
	assert.True(t, zeroDetails.IsZero(), "Zero value struct should return true")

	// Case 2: Valid struct
	validDetails, _ := assetsource.NewBankDetails("X", "Y")
	assert.False(t, validDetails.IsZero(), "Valid struct should return false")
}
