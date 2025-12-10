package wallet

import (
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/assetsource"

	"github.com/google/uuid"
)

type ID uuid.UUID

type Wallet struct {
	id           ID
	name         string
	isStrictMode bool
	currency     valueobject.Currency

	allocations []*Allocation
}

func NewWallet(
	name string,
	currency valueobject.Currency,
	isStrictMode bool,
	allocations []*Allocation,
) (*Wallet, error) {
	if name == "" {
		return nil, errors.New("wallet name cannot be empty")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	if len(allocations) == 0 {
		return nil, errors.New("wallet is not belong to assert source")
	}

	return &Wallet{
		id:           ID(id),
		name:         name,
		currency:     currency,
		isStrictMode: isStrictMode,
		allocations:  allocations,
	}, nil
}

func UnmarshalWalletFromDB(
	id ID,
	name string,
	currency valueobject.Currency,
	isStrictMode bool,
	allocations []*Allocation,
) (*Wallet, error) {
	if id == ID(uuid.Nil) {
		return nil, errors.New("cannot load wallet with nil ID")
	}

	if name == "" {
		return nil, errors.New("data corruption: wallet name is empty")
	}

	if currency.IsZero() {
		return nil, errors.New("data corruption: wallet currency is empty")
	}

	return &Wallet{
		id:           id,
		name:         name,
		currency:     currency,
		isStrictMode: isStrictMode,
		allocations:  allocations,
	}, nil
}

// --- GETTERS (Crucial for other layers to read data) ---
func (w *Wallet) ID() ID                         { return w.id }
func (w *Wallet) Name() string                   { return w.name }
func (w *Wallet) IsStrictMode() bool             { return w.isStrictMode }
func (w *Wallet) Currency() valueobject.Currency { return w.currency }
func (w *Wallet) Allocations() []*Allocation     { return w.allocations }

// --- DOMAIN BEHAVIOR ---
func (w *Wallet) TotalBalance() (valueobject.Money, error) {
	total, err := valueobject.NewMoney(0, w.currency)
	if err != nil {
		return valueobject.Money{}, fmt.Errorf("fail to calculate wallet(ID: %s)'s total balance: %w", uuid.UUID(w.id).String(), err)
	}

	for _, a := range w.allocations {
		total, err = total.Add(a.amount)
		if err != nil {
			return valueobject.Money{}, fmt.Errorf("fail to calculate wallet(ID: %s)'s total balance: %w", uuid.UUID(w.id).String(), err)
		}
	}

	return total, nil
}

func (w *Wallet) TopUp(assetSourceID assetsource.ID, amount valueobject.Money) error {
	if amount.IsZero() {
		return errors.New("top-up amount must be positive")
	}

	if amount.Currency() != w.currency {
		return fmt.Errorf("wallet currency is %s but top-up amount is %s", w.currency.Code(), amount.Currency().Code())
	}

	// Find existing allocation and update it
	for _, alloc := range w.allocations {
		if alloc.assetSourceID == assetSourceID {
			newAmount, err := alloc.amount.Add(amount)
			if err != nil {
				return err
			}

			alloc.amount = newAmount
			return nil
		}
	}

	// returning error if not found source
	return fmt.Errorf("asset source %s not found in this wallet", assetSourceID)
}

func (w *Wallet) Withdraw(assetSourceID assetsource.ID, amount valueobject.Money) error {
	if amount.Currency() != w.currency {
		return fmt.Errorf("wallet currency is %s but withdraw amount is %s", w.currency, amount.Currency())
	}

	for _, alloc := range w.allocations {
		if alloc.assetSourceID == assetSourceID {
			// Check if enough balance
			// Assuming Money.Subtract returns error if result < 0
			newAmount, err := alloc.amount.Subtract(amount)
			if err != nil {
				return fmt.Errorf("insufficient funds in asset source %s: %w", assetSourceID, err)
			}

			alloc.amount = newAmount
			return nil
		}
	}

	// returning error if not found source
	return fmt.Errorf("asset source %s not found in this wallet", assetSourceID)
}
