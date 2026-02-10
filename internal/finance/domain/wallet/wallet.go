package wallet

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

var (
	ErrFundProviderAlreadyRegistered = errors.New("fund provider already registered")
)

type Wallet struct {
	id      uuid.UUID
	balance valueobject.Money

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
		providerManager: &ProviderManager{
			providers: make(map[uuid.UUID]ProviderAllocation),
		},
	}, nil
}

func UnmarshalWalletFromDatabase(
	id uuid.UUID,
	balance valueobject.Money,
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
		providerManager: providerManager,
	}, nil
}

func (w *Wallet) Balance() valueobject.Money        { return w.balance }
func (w *Wallet) ProviderManager() *ProviderManager { return w.providerManager }

func (w *Wallet) AddFundProvider(
	provider *FundProvider,
	allocated valueobject.Money,
) error {
	if !w.isCurrencyValidForAllocation(provider, allocated) {
		return ErrCurrencyMismatch
	}

	err := w.providerManager.AddAndAllocate(provider, allocated)
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

func (w *Wallet) isCurrencyValidForAllocation(fundProvider *FundProvider, allocated valueobject.Money) bool {
	return fundProvider.balance.Currency().Equal(w.walletCurrency()) && allocated.Currency().Equal(w.walletCurrency())
}
