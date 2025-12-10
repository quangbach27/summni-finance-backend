package assetsource_test

import (
	"errors"
	"testing"

	// Must use the full path to access the package's public types
	commons_errors "sumni-finance-backend/internal/common/errors"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/assetsource"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Fixtures ---
var (
	testOwnerID = uuid.New()
	nilOwnerID  = uuid.Nil
	usd, _      = valueobject.NewCurrency("USD")
)

// --- Test NewBankAssetSource Factory ---

func TestNewBankAssetSource_Success(t *testing.T) {
	asset, err := assetsource.NewBankAssetSource(
		testOwnerID,
		1000,
		usd,
		"TestBank",
		"1234567890",
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, asset.ID(), "Asset ID should be generated")
	assert.Equal(t, assetsource.BankType, asset.Type(), "Type must be Bank")
	require.NotNil(t, asset.BankDetails(), "BankDetails must not be nil")
	assert.Equal(t, "TestBank", asset.BankDetails().BankName())
}

func TestNewBankAssetSource_ValidationFailure_BaseAndBankDetails(t *testing.T) {
	// 1. Inputs designed to fail multiple layers of validation
	asset, err := assetsource.NewBankAssetSource(
		nilOwnerID, // Fails base asset validation (ownerID)
		-100,       // Fails base asset validation (NewMoney amount)
		usd,
		"", // Fails bank details validation (bankName)
		"",
	)

	require.Error(t, err)
	assert.Nil(t, asset)

	// 2. Assert the error is the merged ValidationErrors type
	var validationErrs *commons_errors.ValidationErrors
	require.True(t, errors.As(err, &validationErrs), "Error must be of type *ValidationErrors")

	// 3. Verify the count of merged errors (Expecting 3: ownerID, balance, bankDetails)
	assert.Len(t, validationErrs.Errors, 3, "Should have 3 accumulated validation errors")

	// Helper to check for specific error messages based on the client-facing field name
	checkError := func(field string, expectedMsg string) {
		found := false
		for _, e := range validationErrs.Errors {
			if e.Field == field && assert.Contains(t, e.Message, expectedMsg) {
				found = true
				break
			}
		}
		assert.True(t, found, "Did not find error for field: "+field)
	}

	// 4. Check for errors from newBaseAssetSource
	checkError("ownerID", "ownerID is required")
	checkError("balance", "cannot be negative") // Message from NewMoney
}

// --- Test NewCashAssetSource Factory ---

func TestNewCashAssetSource_Success(t *testing.T) {
	asset, err := assetsource.NewCashAssetSource(
		testOwnerID,
		500,
		usd,
	)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, asset.ID())
	assert.Equal(t, assetsource.CashType, asset.Type(), "Type must be Cash")
	assert.Nil(t, asset.BankDetails(), "Cash asset must have nil BankDetails")
}

func TestNewCashAssetSource_ValidationFailure(t *testing.T) {
	asset, err := assetsource.NewCashAssetSource(
		nilOwnerID,             // Fails base asset validation (ownerID)
		-50,                    // Fails base asset validation (amount)
		valueobject.Currency{}, // Fails base asset validation (currency)
	)

	require.Error(t, err)
	assert.Nil(t, asset)

	// Cash factory only relies on newBaseAssetSource, so we test its accumulated errors
	var validationErrs *commons_errors.ValidationErrors
	require.True(t, errors.As(err, &validationErrs))

	// Expecting 3 errors: ownerID, amount, currency
	assert.Len(t, validationErrs.Errors, 3)

	checkError := func(field string, expectedMsg string) {
		found := false
		for _, e := range validationErrs.Errors {
			if e.Field == field && assert.Contains(t, e.Message, expectedMsg) {
				found = true
				break
			}
		}
		assert.True(t, found, "Did not find error for field: "+field)
	}

	checkError("ownerID", "ownerID is required")
	checkError("balance", "cannot be negative")
	checkError("currency", "currency is required")
}

func TestNewSourceTypeFromStr_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType assetsource.SourceType
	}{
		{
			name:     "Success: Exact match (CASH)",
			input:    "CASH",
			wantType: assetsource.CashType,
		},
		{
			name:     "Success: Exact match (BANK)",
			input:    "BANK",
			wantType: assetsource.BankType,
		},
		{
			name:     "Success: Lowercase input (cash)",
			input:    "cash",
			wantType: assetsource.CashType,
		},
		{
			name:     "Success: Mixed case input (BaNk)",
			input:    "BaNk",
			wantType: assetsource.BankType,
		},
		{
			name:     "Success: Input with leading/trailing spaces",
			input:    " BANK ",
			wantType: assetsource.BankType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assetsource.NewSourceTypeFromStr(tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.wantType, got)
			assert.Equal(t, tt.wantType.Code(), got.Code())
			assert.False(t, got.IsZero())
		})
	}
}

func TestNewSourceTypeFromStr_Failure(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Failure: Unknown type",
			input: "CRYPTO",
		},
		{
			name:  "Failure: Empty string",
			input: "",
		},
		{
			name:  "Failure: Whitespace only",
			input: "  ",
		},
		{
			name:  "Failure: Typo",
			input: "CASHH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assetsource.NewSourceTypeFromStr(tt.input)

			require.Error(t, err)
			assert.True(t, got.IsZero(), "Should return zero value on error")
			assert.Contains(t, err.Error(), "unknow asset source type:") // Check for expected error text
		})
	}
}
