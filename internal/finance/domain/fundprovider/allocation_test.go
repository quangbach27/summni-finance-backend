package fundprovider_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	existWalletID = uuid.New()
	baseCurrency  = valueobject.USD
)

func TestNewAllocation(t *testing.T) {
	testCases := []struct {
		name            string
		currency        valueobject.Currency
		entries         []fundprovider.AllocationEntry
		hasError        bool
		expectedTotal   int64
		expectedEntries int
	}{
		{
			name:     "should fail when currency is zero",
			currency: valueobject.Currency{},
			entries:  []fundprovider.AllocationEntry{},
			hasError: true,
		},
		{
			name:            "should success with empty entries",
			currency:        baseCurrency,
			entries:         []fundprovider.AllocationEntry{},
			hasError:        false,
			expectedTotal:   0,
			expectedEntries: 0,
		},
		{
			name:     "should success with single entry",
			currency: baseCurrency,
			entries: []fundprovider.AllocationEntry{
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 100, baseCurrency),
				},
			},
			hasError:        false,
			expectedTotal:   100,
			expectedEntries: 1,
		},
		{
			name:     "should success with multiple entries",
			currency: baseCurrency,
			entries: []fundprovider.AllocationEntry{
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 100, baseCurrency),
				},
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 200, baseCurrency),
				},
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 300, baseCurrency),
				},
			},
			hasError:        false,
			expectedTotal:   600,
			expectedEntries: 3,
		},
		{
			name:     "should fail when entry has nil wallet ID",
			currency: baseCurrency,
			entries: []fundprovider.AllocationEntry{
				{
					WalletID: uuid.Nil,
					Amount:   assertNewMoney(t, 100, baseCurrency),
				},
			},
			hasError: true,
		},
		{
			name:     "should fail when entry has empty amount",
			currency: baseCurrency,
			entries: []fundprovider.AllocationEntry{
				{
					WalletID: uuid.New(),
					Amount:   valueobject.Money{},
				},
			},
			hasError: true,
		},
		{
			name:     "should fail when entry amount currency mismatch",
			currency: baseCurrency,
			entries: []fundprovider.AllocationEntry{
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 100, valueobject.VND),
				},
			},
			hasError: true,
		},
		{
			name:     "should success with zero amount entry",
			currency: baseCurrency,
			entries: []fundprovider.AllocationEntry{
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 0, baseCurrency),
				},
			},
			hasError:        false,
			expectedTotal:   0,
			expectedEntries: 1,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When create new allocation
			allocation, err := fundprovider.NewAllocation(tt.currency, tt.entries...)

			if tt.hasError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTotal, allocation.TotalAllocated().Amount())
			assert.Equal(t, tt.currency, allocation.TotalAllocated().Currency())
			assert.Len(t, allocation.Entries(), tt.expectedEntries)

			// Verify each entry exists in the allocation
			for _, entry := range tt.entries {
				storedEntry, exist := allocation.EntryOf(entry.WalletID)
				assert.True(t, exist)
				assert.Equal(t, entry.WalletID, storedEntry.WalletID)
				assert.Equal(t, entry.Amount, storedEntry.Amount)
			}
		})
	}
}

func TestAllocation_Allocate(t *testing.T) {
	testCases := []struct {
		name          string
		inputWalletID uuid.UUID
		inputMoney    valueobject.Money
		hasError      bool
		expectedError error
	}{
		{
			name:          "should fail when wallets is already allocated",
			inputWalletID: existWalletID,
			inputMoney:    assertNewMoney(t, 10, baseCurrency),
			hasError:      true,
			expectedError: fundprovider.ErrWalletAlreadyAllocated,
		},
		{
			name:          "should fail when money allocations is mismatch currency",
			inputWalletID: uuid.New(),
			inputMoney:    assertNewMoney(t, 10, valueobject.VND),
			hasError:      true,
		},
		{
			name:          "should fail when money allocations is empty",
			inputWalletID: uuid.New(),
			inputMoney:    valueobject.Money{},
			hasError:      true,
			expectedError: fundprovider.ErrInvalidAllocationEntry,
		},
		{
			name:          "should fail when entry has nil id",
			inputWalletID: uuid.Nil,
			inputMoney:    assertNewMoney(t, 10, valueobject.USD),
			hasError:      true,
			expectedError: fundprovider.ErrInvalidAllocationEntry,
		},
		{
			name:          "should success when money allocation is zero",
			inputWalletID: uuid.New(),
			inputMoney:    assertNewMoney(t, 0, valueobject.USD),
			hasError:      false,
		},
		{
			name:          "should allocate success",
			inputWalletID: uuid.New(),
			inputMoney:    assertNewMoney(t, 10, valueobject.USD),
			hasError:      false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given exist wallet allocations
			allocationEntries := []fundprovider.AllocationEntry{
				{
					WalletID: existWalletID,
					Amount:   assertNewMoney(t, 100, baseCurrency),
				},
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 100, baseCurrency),
				},
			}

			// Given wallet Allocation
			walletAllocation, err := fundprovider.NewAllocation(baseCurrency, allocationEntries...)
			assert.NoError(t, err)

			// When perform allocate
			newWalletAllocation, err := walletAllocation.Allocate(tt.inputWalletID, tt.inputMoney)

			if tt.hasError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}

				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, newWalletAllocation)

			walletEntry, exist := newWalletAllocation.EntryOf(tt.inputWalletID)
			assert.True(t, exist)
			assert.Equal(t, tt.inputWalletID, walletEntry.WalletID)
			assert.Equal(t, tt.inputMoney, walletEntry.Amount)
		})
	}
}

func TestAllocation_IncreaseAllocation(t *testing.T) {
	testCases := []struct {
		name           string
		inputWalletID  uuid.UUID
		inputMoney     valueobject.Money
		initialAmount  int64
		expectedAmount int64
		hasError       bool
		expectedError  error
	}{
		{
			name:          "should fail when wallet is not allocated",
			inputWalletID: uuid.New(),
			inputMoney:    assertNewMoney(t, 10, baseCurrency),
			hasError:      true,
			expectedError: fundprovider.ErrWalletNotAllocated,
		},
		{
			name:          "should fail when money is mismatch currency",
			inputWalletID: existWalletID,
			inputMoney:    assertNewMoney(t, 10, valueobject.VND),
			initialAmount: 100,
			hasError:      true,
		},
		{
			name:          "should fail when money is empty",
			inputWalletID: existWalletID,
			inputMoney:    valueobject.Money{},
			initialAmount: 100,
			hasError:      true,
		},
		{
			name:           "should success when increase by zero",
			inputWalletID:  existWalletID,
			inputMoney:     assertNewMoney(t, 0, baseCurrency),
			initialAmount:  100,
			expectedAmount: 100,
			hasError:       false,
		},
		{
			name:           "should increase allocation success",
			inputWalletID:  existWalletID,
			inputMoney:     assertNewMoney(t, 50, baseCurrency),
			initialAmount:  100,
			expectedAmount: 150,
			hasError:       false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given exist wallet allocations
			allocationEntries := []fundprovider.AllocationEntry{
				{
					WalletID: existWalletID,
					Amount:   assertNewMoney(t, tt.initialAmount, baseCurrency),
				},
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 100, baseCurrency),
				},
			}

			// Given wallet Allocation
			walletAllocation, err := fundprovider.NewAllocation(baseCurrency, allocationEntries...)
			assert.NoError(t, err)

			// When perform increase allocation
			newWalletAllocation, err := walletAllocation.IncreaseAllocation(tt.inputWalletID, tt.inputMoney)

			if tt.hasError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}

				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, newWalletAllocation)

			walletEntry, exist := newWalletAllocation.EntryOf(tt.inputWalletID)
			assert.True(t, exist)
			assert.Equal(t, tt.inputWalletID, walletEntry.WalletID)
			assert.Equal(t, tt.expectedAmount, walletEntry.Amount.Amount())
		})
	}
}

func TestAllocation_DecreaseAllocation(t *testing.T) {
	testCases := []struct {
		name           string
		inputWalletID  uuid.UUID
		inputMoney     valueobject.Money
		initialAmount  int64
		expectedAmount int64
		hasError       bool
		expectedError  error
	}{
		{
			name:          "should fail when wallet is not allocated",
			inputWalletID: uuid.New(),
			inputMoney:    assertNewMoney(t, 10, baseCurrency),
			hasError:      true,
			expectedError: fundprovider.ErrWalletNotAllocated,
		},
		{
			name:          "should fail when money is mismatch currency",
			inputWalletID: existWalletID,
			inputMoney:    assertNewMoney(t, 10, valueobject.VND),
			initialAmount: 100,
			hasError:      true,
		},
		{
			name:          "should fail when money is empty",
			inputWalletID: existWalletID,
			inputMoney:    valueobject.Money{},
			initialAmount: 100,
			hasError:      true,
		},
		{
			name:          "should fail when decrease amount exceeds allocation",
			inputWalletID: existWalletID,
			inputMoney:    assertNewMoney(t, 150, baseCurrency),
			initialAmount: 100,
			hasError:      true,
			expectedError: fundprovider.ErrAllocationLimitExceeded,
		},
		{
			name:           "should success when decrease by zero",
			inputWalletID:  existWalletID,
			inputMoney:     assertNewMoney(t, 0, baseCurrency),
			initialAmount:  100,
			expectedAmount: 100,
			hasError:       false,
		},
		{
			name:           "should decrease allocation success",
			inputWalletID:  existWalletID,
			inputMoney:     assertNewMoney(t, 50, baseCurrency),
			initialAmount:  100,
			expectedAmount: 50,
			hasError:       false,
		},
		{
			name:           "should success when decrease to zero",
			inputWalletID:  existWalletID,
			inputMoney:     assertNewMoney(t, 100, baseCurrency),
			initialAmount:  100,
			expectedAmount: 0,
			hasError:       false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Given exist wallet allocations
			allocationEntries := []fundprovider.AllocationEntry{
				{
					WalletID: existWalletID,
					Amount:   assertNewMoney(t, tt.initialAmount, baseCurrency),
				},
				{
					WalletID: uuid.New(),
					Amount:   assertNewMoney(t, 100, baseCurrency),
				},
			}

			// Given wallet Allocation
			walletAllocation, err := fundprovider.NewAllocation(baseCurrency, allocationEntries...)
			assert.NoError(t, err)

			// When perform decrease allocation
			newWalletAllocation, err := walletAllocation.DecreaseAllocation(tt.inputWalletID, tt.inputMoney)

			if tt.hasError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}

				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, newWalletAllocation)

			walletEntry, exist := newWalletAllocation.EntryOf(tt.inputWalletID)
			assert.True(t, exist)
			assert.Equal(t, tt.inputWalletID, walletEntry.WalletID)
			assert.Equal(t, tt.expectedAmount, walletEntry.Amount.Amount())
		})
	}
}

func TestAllocation_IsImmutable(t *testing.T) {
	t.Run("should Allocation immutable when allocate new wallet", func(t *testing.T) {
		t.Parallel()

		baseAllocationEntry := []fundprovider.AllocationEntry{
			{
				WalletID: uuid.New(),
				Amount:   assertNewMoney(t, 1000, baseCurrency),
			},
		}

		original, err := fundprovider.NewAllocation(baseCurrency, baseAllocationEntry...)
		assert.NoError(t, err)

		_, err = original.Allocate(
			uuid.New(),
			assertNewMoney(t, 100, valueobject.USD),
		)
		assert.NoError(t, err)
		assertImmutable(t, original, baseCurrency, baseAllocationEntry)
	})

	t.Run("should Allocation immutable when increase allocation", func(t *testing.T) {
		t.Parallel()

		walletID := uuid.New()
		baseAllocationEntry := []fundprovider.AllocationEntry{
			{
				WalletID: walletID,
				Amount:   assertNewMoney(t, 1000, baseCurrency),
			},
		}

		original, err := fundprovider.NewAllocation(baseCurrency, baseAllocationEntry...)
		assert.NoError(t, err)

		_, err = original.IncreaseAllocation(walletID, assertNewMoney(t, 100, valueobject.USD))
		assert.NoError(t, err)
		assertImmutable(t, original, baseCurrency, baseAllocationEntry)
	})

	t.Run("should Allocation immutable when decrease allocation", func(t *testing.T) {
		t.Parallel()

		walletID := uuid.New()
		baseAllocationEntry := []fundprovider.AllocationEntry{
			{
				WalletID: walletID,
				Amount:   assertNewMoney(t, 1000, baseCurrency),
			},
		}

		original, err := fundprovider.NewAllocation(baseCurrency, baseAllocationEntry...)
		assert.NoError(t, err)

		_, err = original.DecreaseAllocation(walletID, assertNewMoney(t, 100, valueobject.USD))
		assert.NoError(t, err)
		assertImmutable(t, original, baseCurrency, baseAllocationEntry)
	})
}

func assertNewMoney(t *testing.T, amount int64, currency valueobject.Currency) valueobject.Money {
	t.Helper()

	money, err := valueobject.NewMoney(amount, currency)

	assert.NoError(t, err)
	assert.Equal(t, money.Amount(), amount)
	assert.Equal(t, money.Currency(), currency)

	return money
}

func assertImmutable(
	t *testing.T,
	allocation fundprovider.Allocation,
	baseCurrency valueobject.Currency,
	baseAllocationEntry []fundprovider.AllocationEntry,
) {
	t.Helper()

	assert.Equal(t, allocation.TotalAllocated().Currency(), baseCurrency)
	assert.Len(t, allocation.Entries(), len(baseAllocationEntry))
	assert.Equal(t, allocation.Entries(), baseAllocationEntry)
}
