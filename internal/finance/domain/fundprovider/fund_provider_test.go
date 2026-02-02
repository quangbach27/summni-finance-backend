package fundprovider_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	baseCurrency valueobject.Currency = valueobject.USD
)

func TestFundProvider_TopUp(t *testing.T) {
	testCases := []struct {
		name string

		inputBalance     valueobject.Money
		inputTopUpAmount valueobject.Money

		hasErr bool

		expectedBalance valueobject.Money
		expectedErr     error
	}{
		{
			name:             "successful",
			inputBalance:     assertNewMoney(t, 100, baseCurrency),
			inputTopUpAmount: assertNewMoney(t, 10, baseCurrency),
			hasErr:           false,
			expectedBalance:  assertNewMoney(t, 110, baseCurrency),
		},
		{
			name:             "should failed when currency mismatch",
			inputBalance:     assertNewMoney(t, 100, baseCurrency),
			inputTopUpAmount: assertNewMoney(t, 10, valueobject.VND),
			hasErr:           true,
		},
		{
			name:             "should success when top up amount is zero",
			inputBalance:     assertNewMoney(t, 100, baseCurrency),
			inputTopUpAmount: assertNewMoney(t, 0, baseCurrency),
			hasErr:           true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Given
			fundProvider, err := fundprovider.NewFundProvider(
				tt.inputBalance,
				fundprovider.BankProviderType,
				fundprovider.ProviderDetailsOptions{
					BankName:      "Techcombank",
					AccountOwner:  "BUI QUANG BACH",
					AccountNumber: "7777777316",
				},
			)
			require.NoError(t, err)

			err = fundProvider.TopUp(tt.inputTopUpAmount)

			// Then
			if tt.hasErr {
				require.Error(t, err)
				if tt.expectedErr != nil {
					assert.ErrorIs(t, err, tt.expectedErr)
				}

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedBalance, fundProvider.Balance())
		})
	}

	t.Run("should failed when currency mismatch", func(t *testing.T) {

	})

	t.Run("should failed when money is zero", func(t *testing.T) {

	})
}

func assertNewMoney(t *testing.T, amount int64, currency valueobject.Currency) valueobject.Money {
	t.Helper()

	money, err := valueobject.NewMoney(amount, currency)
	require.NoError(t, err)

	return money
}
