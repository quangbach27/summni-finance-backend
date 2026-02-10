package wallet

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"

	"github.com/google/uuid"
)

var (
	ErrCurrencyMismatch      = errors.New("currency mismatch")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrInsufficientAvailable = errors.New("insufficient available amount")
)

var (
	ErrFundProviderAlreadyRegistered = errors.New("fund provider already registered")
)

type Wallet struct {
	id      uuid.UUID
	balance valueobject.Money
	version int32

	providerManager *ProviderManager
}

func NewWallet(currency valueobject.Currency) (*Wallet, error) {
	if currency.IsZero() {
		return nil, errors.New("currency is required")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	balance, err := valueobject.NewMoney(0, currency)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		id:      id,
		balance: balance,
		version: 0,
		providerManager: &ProviderManager{
			providers: make(map[uuid.UUID]ProviderAllocation),
		},
	}, nil
}

func UnmarshalWalletFromDatabase(
	id uuid.UUID,
	balance valueobject.Money,
	version int32,
	providerAllocations []ProviderAllocation,
) (*Wallet, error) {
	v := validator.New()

	v.Check(id != uuid.Nil, "id", "id is required")
	v.Check(!balance.IsZero(), "balance", "balance is required")
	v.Check(balance.Amount() >= 0, "balance", "balance must be positive")

	if err := v.Err(); err != nil {
		return nil, err
	}

	providerManager, err := NewProviderManager(providerAllocations)
	if err != nil {
		return nil, err
	}

	totalAllocated, err := providerManager.CalculateTotalProviderAllocated()
	if err != nil {
		return nil, err
	}

	if !totalAllocated.Equal(balance) {
		return nil, errors.New("total allocated does not match with wallet balance")
	}

	return &Wallet{
		id:              id,
		balance:         balance,
		version:         version,
		providerManager: providerManager,
	}, nil
}

func (w *Wallet) ID() uuid.UUID                     { return w.id }
func (w *Wallet) Balance() valueobject.Money        { return w.balance }
func (w *Wallet) ProviderManager() *ProviderManager { return w.providerManager }
func (w *Wallet) Version() int32                    { return w.version }

func (w *Wallet) AddFundProvider(
	fundProvider *fundprovider.FundProvider,
	allocated valueobject.Money,
) error {
	if fundProvider == nil || allocated.IsZero() {
		return errors.New("FundProvider or allocated is required")
	}

	if !w.isCurrencyValidForAllocation(fundProvider, allocated) {
		return ErrCurrencyMismatch
	}

	err := w.providerManager.AddAndAllocate(fundProvider, allocated)
	if err != nil {
		return err
	}

	newBalance, err := w.balance.Add(allocated)
	if err != nil {
		return err
	}

	w.balance = newBalance

	return nil
}

func (w *Wallet) walletCurrency() valueobject.Currency { return w.balance.Currency() }

func (w *Wallet) isCurrencyValidForAllocation(fundProvider *fundprovider.FundProvider, allocated valueobject.Money) bool {
	return fundProvider.Balance().Currency().Equal(w.walletCurrency()) && allocated.Currency().Equal(w.walletCurrency())
}
