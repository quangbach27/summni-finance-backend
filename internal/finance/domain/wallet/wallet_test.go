package wallet_test

import (
	"sumni-finance-backend/internal/finance/domain/wallet"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWallet(t *testing.T) {
	testCases := []struct {
		name         string
		currencyCode string
		hasErr       bool
	}{
		{
			name:         "cannot init wallet when currency code is empty",
			currencyCode: "",
			hasErr:       true,
		},
		{
			name:         "cannot init wallet when currency code is invalid",
			currencyCode: "INVALID",
			hasErr:       true,
		},
		{
			name:         "can init wallet success",
			currencyCode: "USD",
			hasErr:       false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			wallet, err := wallet.NewWallet(tt.currencyCode)

			if tt.hasErr {
				require.Error(t, err)
				assert.Nil(t, wallet)
			} else {
				require.NoError(t, err)

				var expectedBalance int64 = 0
				assert.Equal(t, expectedBalance, wallet.Balance().Amount())
				assert.Equal(t, tt.currencyCode, wallet.Currency().Code())
			}
		})
	}
}
func TestUnmarshalWalletFromDatabase(t *testing.T) {
	testCases := []struct {
		name          string
		id            uuid.UUID
		balanceAmount int64
		currencyCode  string
		hasErr        bool
	}{
		{
			name:   "cannot init wallet when id is empty",
			id:     uuid.UUID{},
			hasErr: true,
		},
		{
			name:          "cannot init wallet when balance is negative",
			id:            uuid.New(),
			balanceAmount: -10,
			currencyCode:  "USD",
			hasErr:        true,
		},
		{
			name:          "can init wallet when balance is zero",
			id:            uuid.New(),
			balanceAmount: 0,
			currencyCode:  "USD",
			hasErr:        false,
		},
		{
			name:          "cannot init wallet when currency code is empty",
			id:            uuid.New(),
			balanceAmount: 0,
			currencyCode:  "",
			hasErr:        true,
		},
		{
			name:          "cannot init wallet when currency code is not valid",
			id:            uuid.New(),
			balanceAmount: 0,
			currencyCode:  "INVALID",
			hasErr:        true,
		},
		{
			name:          "can init wallet success",
			id:            uuid.New(),
			balanceAmount: 10,
			currencyCode:  "USD",
			hasErr:        false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var version int32 = 1
			wallet, err := wallet.UnmarshalWalletFromDatabase(
				tt.id,
				tt.balanceAmount,
				tt.currencyCode,
				version,
			)

			if tt.hasErr {
				require.Error(t, err)
				assert.Nil(t, wallet)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tt.id, wallet.ID())
				assert.Equal(t, tt.balanceAmount, wallet.Balance().Amount())
				assert.Equal(t, tt.currencyCode, wallet.Currency().Code())
				assert.Equal(t, version, wallet.Version())
			}
		})
	}
}
