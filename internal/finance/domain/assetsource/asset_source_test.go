package assetsource_test

import (
	"errors"
	"testing"

	// Must use the full path to access the package's public types
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/assetsource"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Fixtures ---
var (
	testOwnerID   = uuid.New()
	usd, _        = valueobject.NewCurrency("USD")
	bankName      = "techcombank"
	accountNumber = "7777777317"
)

// --- Test AssetSource Factory ---

func TestAssetSource_NewBankAssetSource(t *testing.T) {
	testCases := []struct {
		name string

		inputOwnerID       uuid.UUID
		inputInitAmount    int64
		inputCurrency      valueobject.Currency
		inputBankName      string
		inputAccountNumner string

		expectedErrorFields []string // Used for asserting specific field failures
		wantErr             bool
	}{
		{
			name:               "Success: Valid Inputs",
			inputOwnerID:       testOwnerID,
			inputInitAmount:    1000,
			inputCurrency:      usd,
			inputBankName:      bankName,
			inputAccountNumner: accountNumber,
			wantErr:            false,
		},
		{
			name:                "Fail: Missing OwnerID",
			inputOwnerID:        uuid.UUID{},
			inputInitAmount:     1000,
			inputCurrency:       usd,
			inputBankName:       bankName,
			inputAccountNumner:  accountNumber,
			expectedErrorFields: []string{"ownerID"},
			wantErr:             true,
		},
		{
			name:                "Fail: Negative Amount",
			inputOwnerID:        testOwnerID,
			inputInitAmount:     -10,
			inputCurrency:       usd,
			inputBankName:       bankName,
			inputAccountNumner:  accountNumber,
			expectedErrorFields: []string{"amount"},
			wantErr:             true,
		},
		{
			name:                "Fail: Missing Currency",
			inputOwnerID:        testOwnerID,
			inputInitAmount:     1000,
			inputCurrency:       valueobject.Currency{},
			inputBankName:       bankName,
			inputAccountNumner:  accountNumber,
			expectedErrorFields: []string{"currency"},
			wantErr:             true,
		},
		{
			name:                "Fail: Missing BankName",
			inputOwnerID:        testOwnerID,
			inputInitAmount:     1000,
			inputCurrency:       usd,
			inputBankName:       "", // Missing local check field
			inputAccountNumner:  accountNumber,
			expectedErrorFields: []string{"bankName"}, // Assuming NewBankDetails combines both errors
			wantErr:             true,
		},
		{
			name:                "Fail: Missing AccountNumber",
			inputOwnerID:        testOwnerID,
			inputInitAmount:     1000,
			inputCurrency:       usd,
			inputBankName:       bankName, // Missing local check field
			inputAccountNumner:  "",
			expectedErrorFields: []string{"accountNumber"}, // Assuming NewBankDetails combines both errors
			wantErr:             true,
		},
		{
			name:                "Fail: Missing all required fields",
			inputOwnerID:        uuid.UUID{}, // Base fail
			inputInitAmount:     -10,
			inputCurrency:       valueobject.Currency{},
			inputBankName:       "",
			inputAccountNumner:  "",
			expectedErrorFields: []string{"ownerID", "amount", "currency", "bankName", "accountNumber"},
			wantErr:             true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := assetsource.NewBankAssetSource(
				tt.inputOwnerID,
				tt.inputInitAmount,
				tt.inputCurrency,
				tt.inputBankName,
				tt.inputAccountNumner,
			)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, asset)

				validationErr := &validator.ErrorList{}

				// Assert the error is our validation type
				assert.True(t, errors.As(err, &validationErr), "Expected error to be a ValidationErrors type")

				// Assert the correct number of errors
				assert.Len(t, validationErr.Errors, len(tt.expectedErrorFields), "Incorrect number of errors collected")

				// Assert the specific error fields are present
				for _, field := range tt.expectedErrorFields {
					assert.Contains(t, validationErr.Errors, field, "Missing expected error field: "+field)
				}
			} else {
				require.NoError(t, err)

				assert.NotEqual(t, uuid.Nil, asset.ID(), "Asset ID should be generated")
				assert.Equal(t, assetsource.BankType, asset.Type(), "Type must be Bank")
				require.NotNil(t, asset.BankDetails(), "BankDetails must not be nil")
				assert.Equal(t, bankName, asset.BankDetails().BankName())
				assert.Equal(t, accountNumber, asset.BankDetails().AccountNumber())
			}
		})
	}
}

func TestAssetSource_NewCashAssetSource(t *testing.T) {
	testCases := []struct {
		name string

		inputOwnerID    uuid.UUID
		inputInitAmount int64
		inputCurrency   valueobject.Currency

		expectedErrorFields []string
		wantErr             bool
	}{
		{
			name:            "Success: Valid Inputs",
			inputOwnerID:    testOwnerID,
			inputInitAmount: 1000,
			inputCurrency:   usd,
			wantErr:         false,
		},
		{
			name:                "Fail: Missing OwnerID",
			inputOwnerID:        uuid.UUID{},
			inputInitAmount:     1000,
			inputCurrency:       usd,
			expectedErrorFields: []string{"ownerID"},
			wantErr:             true,
		},
		{
			name:                "Fail: Negative Amount",
			inputOwnerID:        testOwnerID,
			inputInitAmount:     -10,
			inputCurrency:       usd,
			expectedErrorFields: []string{"amount"},
			wantErr:             true,
		},
		{
			name:                "Fail: Missing Currency",
			inputOwnerID:        testOwnerID,
			inputInitAmount:     1000,
			inputCurrency:       valueobject.Currency{},
			expectedErrorFields: []string{"currency"},
			wantErr:             true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			asset, err := assetsource.NewCashAssetSource(
				tt.inputOwnerID,
				tt.inputInitAmount,
				tt.inputCurrency,
			)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, asset)

				validationErr := &validator.ErrorList{}
				assert.True(t, errors.As(err, &validationErr), "Expected error to be a ValidationErrors type")
				assert.Len(t, validationErr.Errors, len(tt.expectedErrorFields), "Incorrect number of errors collected")

				for _, field := range tt.expectedErrorFields {
					assert.Contains(t, validationErr.Errors, field, "Missing expected error field: "+field)
				}
			} else {
				require.NoError(t, err)

				assert.NotEqual(t, uuid.Nil, asset.ID(), "Asset ID should be generated")
				assert.Equal(t, assetsource.CashType, asset.Type(), "Type must be Cash")
				assert.Nil(t, asset.BankDetails(), "BankDetails must be nil for Cash")
			}
		})
	}
}

// --- Test  AssetSource Getter ---
func TestAssetSource_GetterValues(t *testing.T) {
	// 1. Setup Known Input Values
	expectedBalance, _ := valueobject.NewMoney(100, usd)
	expectedOwnerID := uuid.Must(uuid.NewV7())
	expectedBankType := assetsource.BankType
	expectedBankDetails, _ := assetsource.NewBankDetails(bankName, accountNumber)
	expectedCashType := assetsource.CashType

	// 2. Test Case 1: Bank Asset (with BankDetails)
	t.Run("BankAsset Getters", func(t *testing.T) {
		bankAsset, err := assetsource.NewBankAssetSource(
			expectedOwnerID,
			expectedBalance.Amount(),
			expectedBalance.Currency(),
			expectedBankDetails.BankName(),
			expectedBankDetails.AccountNumber(),
		)

		require.NoError(t, err)

		// Test simple value equality
		assert.Equal(t, expectedOwnerID, bankAsset.OwnerID(), "OwnerID() should match input OwnerID")
		assert.Equal(t, expectedBankType, bankAsset.Type(), "Type() should match input Type")

		// Test Money struct equality (composite value)
		assert.Equal(t, expectedBalance, bankAsset.Balance(), "Balance() should match input Balance struct")

		// Test pointer equality (BankDetails should NOT be nil)
		require.NotNil(t, bankAsset.BankDetails(), "BankDetails() must not be nil for BankType")
		assert.Equal(t, expectedBankDetails, *bankAsset.BankDetails(), "BankDetails() should return the exact same pointer")
	})

	// 3. Test Case 2: Cash Asset (without BankDetails)
	t.Run("CashAsset Getters", func(t *testing.T) {
		cashAsset, err := assetsource.NewCashAssetSource(
			expectedOwnerID,
			expectedBalance.Amount(),
			expectedBalance.Currency(),
		)

		require.NoError(t, err)

		// Test simple value equality
		assert.Equal(t, expectedCashType, cashAsset.Type(), "Type() should match input Type")

		// Test pointer equality (BankDetails MUST be nil)
		assert.Nil(t, cashAsset.BankDetails(), "BankDetails() must be nil for CashType")
	})
}

// --- Test SourceTye Factory ---

func TestSourceType_NewSourceTypeFromStr(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectSourceType assetsource.SourceType
		wantErr          bool
		expectedCode     string
	}{
		{
			name:             "Success: CASH (Uppercase)",
			input:            "CASH",
			expectSourceType: assetsource.CashType,
			expectedCode:     "CASH",
			wantErr:          false,
		},
		{
			name:             "Success: cash (Lowercase)",
			input:            "cash",
			expectSourceType: assetsource.CashType,
			expectedCode:     "CASH",
			wantErr:          false,
		},
		{
			name:             "Success: BANK (Uppercase)",
			input:            "BANK",
			expectSourceType: assetsource.BankType,
			expectedCode:     "BANK",
			wantErr:          false,
		},
		{
			name:             "Success: bank (Trailing Space)",
			input:            " bank ",
			expectSourceType: assetsource.BankType,
			expectedCode:     "BANK",
			wantErr:          false,
		},
		{
			name:             "Failure: Unknown String",
			input:            "unknown",
			expectSourceType: assetsource.SourceType{},
			wantErr:          true,
		},
		{
			name:             "Failure: Empty String",
			input:            "",
			expectSourceType: assetsource.SourceType{},
			wantErr:          true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assetsource.NewSourceTypeFromStr(tt.input)

			if tt.wantErr {
				require.Error(t, err, "Expected an error for input: %s", tt.input)
				assert.Equal(t, tt.expectSourceType, got, "Expected zero value on failure")
			} else {
				require.NoError(t, err, "Did not expect an error for input: %s", tt.input)
				assert.Equal(t, tt.expectSourceType, got, "SourceType object mismatch")
				assert.Equal(t, tt.expectedCode, got.Code(), "Code mismatch")
			}
		})
	}
}

func TestSourceType_Getter(t *testing.T) {
	t.Run("CashType Code", func(t *testing.T) {
		assert.Equal(t, "CASH", assetsource.CashType.Code(), "CashType Code() should return 'CASH'")
	})

	t.Run("BankType Code", func(t *testing.T) {
		assert.Equal(t, "BANK", assetsource.BankType.Code(), "BankType Code() should return 'BANK'")
	})
}

func TestSourceType_IsZero(t *testing.T) {
	// Define the expected zero value for the SourceType
	var zeroSourceType assetsource.SourceType

	t.Run("Non-Zero Constants", func(t *testing.T) {
		assert.False(t, assetsource.CashType.IsZero(), "CashType should not be zero")
		assert.False(t, assetsource.BankType.IsZero(), "BankType should not be zero")
	})

	t.Run("Zero Value", func(t *testing.T) {
		assert.True(t, zeroSourceType.IsZero(), "Zero-initialized SourceType should be zero")
	})
}
