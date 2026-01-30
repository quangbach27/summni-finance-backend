package fundprovider_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	PROVIDER_BALANCE int64 = 1000
)

func TestFundProvider_AllocateToWallet(t *testing.T) {
	testCases := []struct {
		name                      string
		inputAllocation           fundprovider.Allocation
		inputNewAllocatedWalletID uuid.UUID
		inputNewAllocatedAmount   valueobject.Money
		hasError                  bool
		expectedError             error
	}{
		{
			name: "should success",
			inputAllocation: assertNewAllocation(t, fundprovider.AllocationEntry{
				WalletID: existWalletID,
				Amount:   assertNewMoney(t, 100, baseCurrency),
			}),
			inputNewAllocatedWalletID: uuid.New(),
			inputNewAllocatedAmount:   assertNewMoney(t, 100, baseCurrency),
			hasError:                  false,
		},
	}

	for _, tt := range testCases {
		// given fundProvider
		fundProvider := assertNewFundProvider(
			t,
			PROVIDER_BALANCE,
			fundprovider.BankProviderType,
			tt.inputAllocation,
		)

		// when allocate wallet
		err := fundProvider.AllocateToWallet(tt.inputNewAllocatedWalletID, tt.inputNewAllocatedAmount)
		if tt.hasError {
			assert.NoError(t, err)
			return
		}

		entry, exist := fundProvider.Allocation().EntryOf(tt.inputNewAllocatedWalletID)

		assert.True(t, exist)
		assert.Equal(t, entry.WalletID, tt.inputNewAllocatedWalletID)
		assert.Equal(t, entry.Amount, tt.inputNewAllocatedAmount)
	}

	t.Run("should success", func(t *testing.T) {
		// given givenAllocation
		givenAllocation, err := fundprovider.NewAllocation(baseCurrency, fundprovider.AllocationEntry{
			Amount:   assertNewMoney(t, 100, baseCurrency),
			WalletID: uuid.New(),
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, givenAllocation)

		// given fundProvider
		fundprovider := assertNewFundProvider(
			t,
			PROVIDER_BALANCE,
			fundprovider.BankProviderType,
			givenAllocation,
		)

		givenNewWalletID := uuid.New()
		givenNewAllocatedAmount := assertNewMoney(t, 100, baseCurrency)

		// When allocate
		err = fundprovider.AllocateToWallet(givenNewWalletID, givenNewAllocatedAmount)
		assert.NoError(t, err)

		entry, exist := fundprovider.Allocation().EntryOf(givenNewWalletID)

		assert.True(t, exist)
		assert.Equal(t, entry.WalletID, givenNewWalletID)
		assert.Equal(t, entry.Amount, givenNewAllocatedAmount)
	})
}

func assertNewFundProvider(
	t *testing.T,
	balanceAmount int64,
	providerType fundprovider.FundProviderType,
	allocation fundprovider.Allocation,
) *fundprovider.FundProvider {
	t.Helper()

	balance, err := valueobject.NewMoney(balanceAmount, baseCurrency)
	assert.NoError(t, err)

	options := fundprovider.ProviderDetailsOptions{}
	if providerType.IsBank() {
		options = fundprovider.ProviderDetailsOptions{
			BankName:      "Test Bank",
			AccountOwner:  "Test Owner",
			AccountNumber: "123456789",
		}
	} else if providerType.IsCash() {
		options = fundprovider.ProviderDetailsOptions{
			Name: "Test Cash",
		}
	}

	fundProvider, err := fundprovider.NewFundProvider(
		balance,
		providerType,
		allocation,
		options,
	)
	require.NoError(t, err)

	return fundProvider
}

func assertNewAllocation(
	t *testing.T,
	entries ...fundprovider.AllocationEntry,
) fundprovider.Allocation {
	t.Helper()

	fundProvider, err := fundprovider.NewAllocation(baseCurrency, entries...)
	require.NoError(t, err)

	return fundProvider
}
