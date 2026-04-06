package wallet_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/ledger"
	"sumni-finance-backend/internal/finance/domain/wallet"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLedgerManager_OpenAccountingPeriod(t *testing.T) {
	validYearMonth := NewValidYearMonth(t, 4, 2026)
	validCurrency, err := valueobject.NewCurrency("VND")
	require.NoError(t, err)
	validMoney, err := valueobject.NewMoney(1_000_000, validCurrency)
	require.NoError(t, err)

	tests := []struct {
		name           string
		inputYearMonth ledger.YearMonth
		inputBalance   valueobject.Money
		setupManager   func(t *testing.T) *wallet.LedgerManager
		hasErr         bool
		errorContains  string
	}{
		{
			name:           "returns error when yearMonth is zero",
			inputYearMonth: ledger.YearMonth{},
			inputBalance:   validMoney,
			setupManager: func(t *testing.T) *wallet.LedgerManager {
				t.Helper()
				m, err := wallet.NewLedgerManager(nil)
				require.NoError(t, err)
				return m
			},
			hasErr:        true,
			errorContains: "year and month is required",
		},
		{
			name:           "returns error when openBalance currency is missing",
			inputYearMonth: validYearMonth,
			inputBalance:   valueobject.Money{},
			setupManager: func(t *testing.T) *wallet.LedgerManager {
				t.Helper()
				m, err := wallet.NewLedgerManager(nil)
				require.NoError(t, err)
				return m
			},
			hasErr: true,
		},
		{
			name:           "returns error when accounting period already opened",
			inputYearMonth: validYearMonth,
			inputBalance:   validMoney,
			setupManager: func(t *testing.T) *wallet.LedgerManager {
				t.Helper()
				m, err := wallet.NewLedgerManager(nil)
				require.NoError(t, err)
				err = m.OpenAccountingPeriod(validYearMonth, validMoney)
				require.NoError(t, err)
				return m
			},
			hasErr:        true,
			errorContains: "already opened",
		},
		{
			name:           "opens new accounting period successfully",
			inputYearMonth: validYearMonth,
			inputBalance:   validMoney,
			setupManager: func(t *testing.T) *wallet.LedgerManager {
				t.Helper()
				m, err := wallet.NewLedgerManager(nil)
				require.NoError(t, err)
				return m
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			manager := tt.setupManager(t)

			// Act
			err := manager.OpenAccountingPeriod(tt.inputYearMonth, tt.inputBalance)

			// Assert
			if tt.hasErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.ErrorContains(t, err, tt.errorContains)
				}
				return
			}

			require.NoError(t, err)
			ap, exists := manager.FindAccountingPeriod(tt.inputYearMonth)
			require.True(t, exists)
			require.NotNil(t, ap)
		})
	}
}

func NewValidYearMonth(t *testing.T, month, year int) ledger.YearMonth {
	t.Helper()

	yearMonth, err := ledger.NewYearMonth(month, year)
	require.NoError(t, err)

	return yearMonth
}
