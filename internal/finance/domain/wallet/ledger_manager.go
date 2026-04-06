package wallet

import (
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/ledger"
)

type LedgerConfig struct {
	startDate ledger.PeriodStartDay // date of the month
	interval  int32                 // month
}

func NewDefaultLedgerConfig() (LedgerConfig, error) {
	startDate, err := ledger.NewPeriodStartDay(1)
	if err != nil {
		return LedgerConfig{}, nil
	}

	return LedgerConfig{
		startDate: startDate,
		interval:  1,
	}, nil
}

func (lc LedgerConfig) StartDate() ledger.PeriodStartDay { return lc.startDate }
func (lc LedgerConfig) Interval() int32                  { return lc.interval }

type LedgerManager struct {
	config LedgerConfig

	accountPeriods map[ledger.YearMonth]*ledger.AccountingPeriod
}

func NewLedgerManager(accountPeriods []*ledger.AccountingPeriod) (*LedgerManager, error) {
	ledgerConfig, err := NewDefaultLedgerConfig()
	if err != nil {
		return nil, err
	}

	// Initialize map with appropriate capacity
	capacity := len(accountPeriods)
	if capacity == 0 {
		capacity = 1 // Pre-allocate for at least one period
	}

	ledgerManager := &LedgerManager{
		config:         ledgerConfig,
		accountPeriods: make(map[ledger.YearMonth]*ledger.AccountingPeriod, capacity),
	}

	for _, ap := range accountPeriods {
		if ap == nil {
			return nil, errors.New("accounting period can not be nil")
		}

		ledgerManager.accountPeriods[ap.YearMonth()] = ap
	}

	return ledgerManager, nil
}

func (m *LedgerManager) FindAccountingPeriod(yearMonth ledger.YearMonth) (*ledger.AccountingPeriod, bool) {
	ap, exist := m.accountPeriods[yearMonth]
	if !exist || ap == nil {
		return nil, false
	}

	return ap, true
}

func (m *LedgerManager) OpenAccountingPeriod(
	yearMonth ledger.YearMonth,
	openBalance valueobject.Money,
) error {
	if yearMonth.IsZero() {
		return errors.New("open account period: year and month is required")
	}

	if _, exist := m.FindAccountingPeriod(yearMonth); exist {
		return fmt.Errorf("accounting period %s already opened", yearMonth.String())
	}

	newAccountingPeriod, err := ledger.OpenAccountingPeriod(
		yearMonth,
		openBalance,
		m.config.startDate,
		m.config.interval,
	)
	if err != nil {
		return fmt.Errorf("open new period: %w", err)
	}

	m.accountPeriods[newAccountingPeriod.YearMonth()] = newAccountingPeriod
	return nil
}

func (m *LedgerManager) Record(yearMonth ledger.YearMonth, txRecord ledger.TransactionRecord) error {
	ap, exist := m.FindAccountingPeriod(yearMonth)
	if !exist {
		return fmt.Errorf("account period %s not found", yearMonth.String())
	}

	return ap.Record(txRecord)
}
