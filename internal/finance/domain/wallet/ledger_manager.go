package wallet

import (
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/ledger"

	"github.com/google/uuid"
)

type LedgerConfig struct {
	startDate ledger.PeriodStartDay // day
	interval  int32                 // month
}

func (lc LedgerConfig) StartDate() ledger.PeriodStartDay { return lc.startDate }
func (lc LedgerConfig) Interval() int32                  { return lc.interval }

type LedgerManager struct {
	config LedgerConfig

	accountPeriods map[uuid.UUID]*ledger.AccountingPeriod
}

func NewLedgerManager(accountPeriods []*ledger.AccountingPeriod) (*LedgerManager, error) {
	startDay, err := ledger.NewPeriodStartDay(1)
	if err != nil {
		return nil, err
	}

	// Initialize map with appropriate capacity
	capacity := len(accountPeriods)
	if capacity == 0 {
		capacity = 1 // Pre-allocate for at least one period
	}

	ledgerManager := &LedgerManager{
		config: LedgerConfig{
			startDate: startDay,
			interval:  1,
		},
		accountPeriods: make(map[uuid.UUID]*ledger.AccountingPeriod, capacity),
	}

	for _, ap := range accountPeriods {
		if ap == nil {
			return nil, errors.New("accounting period can not be nil")
		}

		ledgerManager.accountPeriods[ap.ID()] = ap
	}

	return ledgerManager, nil
}

func (m *LedgerManager) FindAccountingPeriod(id uuid.UUID) (*ledger.AccountingPeriod, bool) {
	ap, exist := m.accountPeriods[id]
	if !exist || ap == nil {
		return nil, false
	}

	return ap, true
}

func (m *LedgerManager) OpenNewAccountingPeriod(
	yearMonth ledger.YearMonth,
	openBalance valueobject.Money,
) (uuid.UUID, error) {
	if yearMonth.IsZero() {
		return uuid.Nil, errors.New("open account period: year and month is required")
	}

	newAccountingPeriod, err := ledger.OpenAccountingPeriod(
		yearMonth,
		openBalance,
		m.config.startDate,
		m.config.interval,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("open new period: %w", err)
	}

	m.accountPeriods[newAccountingPeriod.ID()] = newAccountingPeriod
	return newAccountingPeriod.ID(), nil
}

func (m *LedgerManager) Record(id uuid.UUID, txRecord ledger.TransactionRecord) error {
	ap, exist := m.FindAccountingPeriod(id)
	if !exist {
		return fmt.Errorf("account period %s not found", id.String())
	}

	return ap.Record(txRecord)
}
