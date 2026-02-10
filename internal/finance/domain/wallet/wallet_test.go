package wallet_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/wallet"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWallet_AddFundProvider(t *testing.T) {
	t.Run("should fail when fundProvider is nil", func(t *testing.T) {
		walletDomain, err := wallet.NewWallet(valueobject.USD)
		require.NoError(t, err)

		err = walletDomain.AddFundProvider(nil, assertNewMoney(t, 100, valueobject.USD))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "FundProvider or allocated is required")
	})

	t.Run("should fail when fundProvider is empty", func(t *testing.T) {
		walletDomain, err := wallet.NewWallet(valueobject.USD)
		require.NoError(t, err)

		err = walletDomain.AddFundProvider(&fundprovider.FundProvider{}, assertNewMoney(t, 100, valueobject.USD))
		require.Error(t, err)
		assert.ErrorIs(t, err, wallet.ErrCurrencyMismatch)
	})

	t.Run("should fail when allocated amount is empty", func(t *testing.T) {
		walletDomain, err := wallet.NewWallet(valueobject.USD)
		require.NoError(t, err)

		err = walletDomain.AddFundProvider(
			&fundprovider.FundProvider{},
			valueobject.Money{},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "FundProvider or allocated is required")
	})

	t.Run("should failed when provider is exist", func(t *testing.T) {
		walletBalance := assertNewMoney(t, 10, valueobject.USD)

		// Assume already has allocated fundProvider for wallet
		newBalanceFundProvider1 := assertNewMoney(t, 100, valueobject.USD)
		fundProvider1 := assertNewFundProvider(t, newBalanceFundProvider1)
		allocated1 := walletBalance
		providerAllocation1, err := wallet.NewProviderAllocation(fundProvider1, allocated1)
		require.NoError(t, err)

		walletDomain, err := wallet.UnmarshalWalletFromDatabase(
			uuid.New(),
			walletBalance,
			[]wallet.ProviderAllocation{providerAllocation1},
		)
		require.NoError(t, err)

		// The FundProvider will be allocated
		newAllocated := assertNewMoney(t, 10, valueobject.USD)

		// When
		err = walletDomain.AddFundProvider(fundProvider1, newAllocated)
		require.Error(t, err)
		assert.ErrorIs(t, err, wallet.ErrFundProviderAlreadyRegistered)
	})

	t.Run("should failed when allocation mismatch currency with wallet", func(t *testing.T) {
		walletBalance := assertNewMoney(t, 10, valueobject.USD)

		// Assume already has allocated fundProvider for wallet
		newBalanceFundProvider1 := assertNewMoney(t, 100, valueobject.USD)
		fundProvider1 := assertNewFundProvider(t, newBalanceFundProvider1)
		allocated1 := walletBalance
		providerAllocation1, err := wallet.NewProviderAllocation(fundProvider1, allocated1)
		require.NoError(t, err)

		walletDomain, err := wallet.UnmarshalWalletFromDatabase(
			uuid.New(),
			walletBalance,
			[]wallet.ProviderAllocation{providerAllocation1},
		)
		require.NoError(t, err)

		// The FundProvider will be allocated
		newFundProviderBalance := assertNewMoney(t, 100, valueobject.USD)
		newFundProvider := assertNewFundProvider(t, newFundProviderBalance)
		newAllocated := assertNewMoney(t, 10, valueobject.VND)

		// When
		err = walletDomain.AddFundProvider(newFundProvider, newAllocated)
		require.Error(t, err)
		assert.ErrorIs(t, err, wallet.ErrCurrencyMismatch)
	})

	t.Run("should failed when provider does not have enough available amount for allocation", func(t *testing.T) {
		walletBalance := assertNewMoney(t, 10, valueobject.USD)

		// Assume already has allocated fundProvider for wallet
		newBalanceFundProvider1 := assertNewMoney(t, 100, valueobject.USD)
		fundProvider1 := assertNewFundProvider(t, newBalanceFundProvider1)
		allocated1 := walletBalance
		providerAllocation1, err := wallet.NewProviderAllocation(fundProvider1, allocated1)
		require.NoError(t, err)

		walletDomain, err := wallet.UnmarshalWalletFromDatabase(
			uuid.New(),
			walletBalance,
			[]wallet.ProviderAllocation{providerAllocation1},
		)
		require.NoError(t, err)

		// The FundProvider will be allocated
		newFundProviderBalance := assertNewMoney(t, 100, valueobject.USD)
		newFundProvider := assertNewFundProvider(t, newFundProviderBalance)
		newAllocated := assertNewMoney(t, 110, valueobject.USD)

		// When
		err = walletDomain.AddFundProvider(newFundProvider, newAllocated)
		require.Error(t, err)
		assert.ErrorIs(t, err, fundprovider.ErrInsufficientAvailable)
	})

	t.Run("should success when provider have enough available amount and provider does not exist", func(t *testing.T) {
		walletBalance := assertNewMoney(t, 10, valueobject.USD)

		// Assume already has allocated fundProvider for wallet
		newBalanceFundProvider1 := assertNewMoney(t, 100, valueobject.USD)
		fundProvider1 := assertNewFundProvider(t, newBalanceFundProvider1)
		allocated1 := walletBalance
		providerAllocation1, err := wallet.NewProviderAllocation(fundProvider1, allocated1)
		require.NoError(t, err)

		walletDomain, err := wallet.UnmarshalWalletFromDatabase(
			uuid.New(),
			walletBalance,
			[]wallet.ProviderAllocation{providerAllocation1},
		)
		require.NoError(t, err)

		// The FundProvider will be allocated
		newFundProviderBalance := assertNewMoney(t, 100, valueobject.USD)
		newFundProvider := assertNewFundProvider(t, newFundProviderBalance)
		newAllocated := assertNewMoney(t, 10, valueobject.USD)

		// When
		err = walletDomain.AddFundProvider(newFundProvider, newAllocated)
		require.NoError(t, err)

		// Then
		// Check wallet balance
		expectedWalletBalance, err := newAllocated.Add(allocated1)
		require.NoError(t, err)

		assert.True(t, walletDomain.Balance().Equal(expectedWalletBalance))

		// Check allocation manager
		// Check total allocated
		totalAllocated, err := walletDomain.ProviderManager().CalculateTotalProviderAllocated()
		require.NoError(t, err)

		assert.True(t, totalAllocated.Equal(expectedWalletBalance))
		assert.True(t, walletDomain.ProviderManager().HasFundProvider(newFundProvider.ID()))

		// Check avaibleBalance of fundProvider
		expectedAvailableBalance, err := newFundProviderBalance.Subtract(newAllocated)
		require.NoError(t, err)
		assert.True(t, walletDomain.ProviderManager().GetFundProvider(newFundProvider.ID()).AvailableAmountForAllocation().Equal(expectedAvailableBalance))
	})
}

func assertNewFundProvider(t *testing.T, balance valueobject.Money) *fundprovider.FundProvider {
	t.Helper()

	fundProvider, err := fundprovider.NewFundProvider(balance)
	require.NoError(t, err, "newFundProvider should not have error")

	return fundProvider
}

func assertNewMoney(t *testing.T, amount int64, currency valueobject.Currency) valueobject.Money {
	t.Helper()

	money, err := valueobject.NewMoney(amount, currency)
	require.NoError(t, err, "newMoney should not have error")

	return money
}
