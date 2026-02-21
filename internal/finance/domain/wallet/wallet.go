package wallet

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"

	"github.com/google/uuid"
)

var (
	ErrCurrencyMismatch              = errors.New("currency mismatch")
	ErrInsufficientBalance           = errors.New("insufficient balance")
	ErrInsufficientAvailable         = errors.New("insufficient available amount")
	ErrFundProviderAlreadyRegistered = errors.New("fund provider already registered")
	ErrFundAllocatedMissing          = errors.New("fund provider for allocation is missing")
	ErrAllocationAmountNegative      = errors.New("allocated amount is negative")
)

// Wallet is the Root Aggregate
// It contain the FundProvider entity
type Wallet struct {
	id      uuid.UUID
	balance valueobject.Money
	version int32

	providerManager *ProviderManager
}

func NewWallet(currencyCode string) (*Wallet, error) {
	currency, err := valueobject.NewCurrency(currencyCode)
	if err != nil {
		return nil, err
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
	balanceAmount int64,
	currencyCode string,
	version int32,
	providerAllocations ...ProviderAllocation,
) (*Wallet, error) {
	v := validator.New()

	v.Check(id != uuid.Nil, "id", "id is required")
	v.Check(balanceAmount >= 0, "balance", "balance must greater or equal than 0")
	v.Required(currencyCode, "currency")

	if err := v.Err(); err != nil {
		return nil, err
	}

	currency, err := valueobject.NewCurrency(currencyCode)
	if err != nil {
		return nil, err
	}

	balance, err := valueobject.NewMoney(balanceAmount, currency)
	if err != nil {
		return nil, err
	}

	providerManager, err := NewProviderManager(providerAllocations)
	if err != nil {
		return nil, err
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
func (w *Wallet) Currency() valueobject.Currency    { return w.balance.Currency() }
func (w *Wallet) ProviderManager() *ProviderManager { return w.providerManager }
func (w *Wallet) Version() int32                    { return w.version }

func (w *Wallet) AllocateFromFundProvider(
	fundProvider *fundprovider.FundProvider,
	allocatedAmount int64,
) error {
	if fundProvider == nil {
		return ErrFundAllocatedMissing
	}

	if allocatedAmount < 0 {
		return ErrAllocationAmountNegative
	}

	if _, exists := w.ProviderManager().FindProvider(fundProvider.ID()); exists {
		return ErrFundProviderAlreadyRegistered
	}

	allocated, err := valueobject.NewMoney(allocatedAmount, w.Currency())
	if err != nil {
		return err
	}

	return w.providerManager.AddFundProviderAndReserve(fundProvider, allocated)
}
