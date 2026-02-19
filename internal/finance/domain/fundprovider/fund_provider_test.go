package fundprovider_test

import (
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFundProvider_NewFundProvider(t *testing.T) {
	testCases := []struct {
		name              string
		initBalanceAmount int64
		currencyCode      string
		hasErr            bool
		expectedErr       string
	}{
		{
			name:              "cannot init fund provider when init balance is negative",
			initBalanceAmount: -100,
			currencyCode:      "USD",
			hasErr:            true,
		},
		{
			name:              "can init fund provider when init balance is zero",
			initBalanceAmount: 0,
			currencyCode:      "USD",
			hasErr:            false,
		},
		{
			name:              "can init fund provider when init balance is positive",
			initBalanceAmount: 100,
			currencyCode:      "USD",
			hasErr:            false,
		},
		{
			name:              "cannot init fund provider when currency code is empty",
			initBalanceAmount: 0,
			currencyCode:      "",
			hasErr:            true,
		},
		{
			name:              "cannot init fund provider when currency code is invalid",
			initBalanceAmount: 0,
			currencyCode:      "INVALID",
			hasErr:            true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fundProvider, err := fundprovider.NewFundProvider(tt.initBalanceAmount, tt.currencyCode)

			if tt.hasErr {
				require.Error(t, err)
				assert.Nil(t, fundProvider)
			} else {
				require.NoError(t, err)

				assert.Equal(t, fundProvider.Balance().Amount(), tt.initBalanceAmount)
				assert.Equal(t, fundProvider.Balance().Currency().Code(), tt.currencyCode)
			}
		})
	}
}

func TestFundProvider_UnmarshallFromDatabase(t *testing.T) {
	testCases := []struct {
		name string

		// Given
		id                       uuid.UUID
		balanceAmount            int64
		unallocatedBalanceAmount int64
		currencyCode             string
		version                  int32

		// Then
		hasErr bool
	}{
		{
			name:   "cannot initialize when id is empty",
			id:     uuid.UUID{},
			hasErr: true,
		},
		{
			name:          "cannot initialize fund provider when balanceAmount is negative",
			id:            uuid.New(),
			balanceAmount: -100,
			hasErr:        true,
		},
		{
			name:                     "can initialize fund provider when balanceAmount is zero",
			id:                       uuid.New(),
			balanceAmount:            0,
			unallocatedBalanceAmount: 0,
			currencyCode:             "USD",
			version:                  0,
			hasErr:                   false,
		},
		{
			name:                     "can initialize fund provider when balanceAmount is positive",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: 0,
			currencyCode:             "USD",
			version:                  0,
			hasErr:                   false,
		},
		{
			name:                     "cannot initialize fund provider when unallocated amount is negative",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: -10,
			currencyCode:             "USD",
			version:                  0,
			hasErr:                   true,
		},
		{
			name:                     "cannot initialize fund provider when unallocated amount exceed balance amount",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: 120,
			currencyCode:             "USD",
			version:                  0,
			hasErr:                   true,
		},
		{
			name:                     "can initialize fund provider when unallocated amount is zero",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: 0,
			currencyCode:             "USD",
			version:                  0,
			hasErr:                   false,
		},
		{
			name:                     "can initialize fund provider when unallocated amount does not excced balance",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: 50,
			currencyCode:             "USD",
			version:                  0,
			hasErr:                   false,
		},
		{
			name:                     "cannot initialize fund provider when currency code is empty",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: 50,
			currencyCode:             "",
			version:                  0,
			hasErr:                   true,
		},
		{
			name:                     "cannot initialize fund provider when currency code is invalid",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: 50,
			currencyCode:             "INVALID",
			version:                  0,
			hasErr:                   true,
		},
		{
			name:                     "can initialize fund provider",
			id:                       uuid.New(),
			balanceAmount:            100,
			unallocatedBalanceAmount: 50,
			currencyCode:             "USD",
			version:                  1,
			hasErr:                   false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
				tt.id,
				tt.balanceAmount,
				tt.unallocatedBalanceAmount,
				tt.currencyCode,
				tt.version,
			)

			if tt.hasErr {
				require.Error(t, err)
				assert.Nil(t, fundProvider)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tt.id, fundProvider.ID())

				assert.Equal(t, tt.balanceAmount, fundProvider.Balance().Amount())
				assert.Equal(t, tt.currencyCode, fundProvider.Balance().Currency().Code())

				assert.Equal(t, tt.unallocatedBalanceAmount, fundProvider.UnallocatedBalance().Amount())
				assert.Equal(t, tt.currencyCode, fundProvider.UnallocatedBalance().Currency().Code())

				assert.Equal(t, tt.version, fundProvider.Version())
			}
		})
	}
}

func TestFundProvider_Allocate(t *testing.T) {
	testCases := []struct {
		name              string
		unallocatedAmount int64
		allocatedAmount   int64
		hasErr            bool
	}{
		{
			name:              "cannot allocate when allocated amount exceed",
			unallocatedAmount: 50,
			allocatedAmount:   100,
			hasErr:            true,
		},
		{
			name:              "cannot allocate when allocated amount is negative",
			unallocatedAmount: 50,
			allocatedAmount:   -10,
			hasErr:            true,
		},
		{
			name:              "can allocate when allocated is zero",
			unallocatedAmount: 50,
			allocatedAmount:   0,
			hasErr:            false,
		},
		{
			name:              "can allocate",
			unallocatedAmount: 50,
			allocatedAmount:   10,
			hasErr:            false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given
			fundProvider, err := fundprovider.NewFundProvider(tt.unallocatedAmount, "USD")
			require.NoError(t, err)
			require.NotNil(t, fundProvider)

			// When
			err = fundProvider.Allocate(tt.allocatedAmount)

			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				expectedUnallocated := tt.unallocatedAmount - tt.allocatedAmount
				assert.Equal(t, expectedUnallocated, fundProvider.UnallocatedBalance().Amount())
			}
		})
	}
}

func TestFundProvider_TopUp(t *testing.T) {
	testCases := []struct {
		name        string
		initBalance int64
		topUpAmount int64
		hasErr      bool
		expectErr   error
	}{
		{
			name:        "cannot top up when amount is negative",
			initBalance: 100,
			topUpAmount: -50,
			hasErr:      true,
			expectErr:   fundprovider.ErrInsufficientAmount,
		},
		{
			name:        "can top up",
			initBalance: 100,
			topUpAmount: 50,
			hasErr:      false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fundProvider, err := fundprovider.NewFundProvider(tt.initBalance, "USD")
			require.NoError(t, err)

			err = fundProvider.TopUp(tt.topUpAmount)

			if tt.hasErr {
				require.Error(t, err)
				assert.ErrorIs(t, tt.expectErr, err)
			} else {
				require.NoError(t, err)

				expectBalance := tt.initBalance + tt.topUpAmount
				assert.Equal(t, expectBalance, fundProvider.Balance().Amount())
			}
		})
	}
}

func TestFundProvider_Withdraw(t *testing.T) {
	testCases := []struct {
		name               string
		balance            int64
		unallocatedBalance int64
		withdrawAmount     int64
		hasErr             bool
		expectErr          error
	}{
		{
			name:               "cannot withdraw when amount is negative",
			balance:            100,
			unallocatedBalance: 50,
			withdrawAmount:     -50,
			hasErr:             true,
			expectErr:          fundprovider.ErrInsufficientAmount,
		},
		{
			name:               "cannot withdraw when amount excceed unallocated balance",
			balance:            100,
			unallocatedBalance: 50,
			withdrawAmount:     60,
			hasErr:             true,
			expectErr:          fundprovider.ErrInsufficientAmount,
		},
		{
			name:               "can withdraw success",
			balance:            100,
			unallocatedBalance: 50,
			withdrawAmount:     30,
			hasErr:             false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			fundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
				uuid.New(),
				tt.balance,
				tt.unallocatedBalance,
				"USD",
				0,
			)
			require.NoError(t, err)

			err = fundProvider.Withdraw(tt.withdrawAmount)

			if tt.hasErr {
				require.Error(t, err)
				assert.ErrorIs(t, tt.expectErr, err)
			} else {
				require.NoError(t, err)

				expectBalance := tt.balance - tt.withdrawAmount
				assert.Equal(t, expectBalance, fundProvider.Balance().Amount())
			}
		})
	}
}
