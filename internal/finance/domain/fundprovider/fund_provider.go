package fundprovider

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

var (
	// Provider-level errors
	ErrInsufficientBalance = errors.New("insufficient fund provider balance")
	ErrCurrencyMismatch    = errors.New("currency mismatch between money and provider")
)

type FundProvider struct {
	id               uuid.UUID
	balance          valueobject.Money
	currency         valueobject.Currency
	fundProviderType FundProviderType
	allocation       Allocation

	providerDetails ProviderDetails
	aggVerion       int32
}

func NewFundProvider(
	balance valueobject.Money,
	fundProviderType FundProviderType,
	allocation Allocation,
	options ProviderDetailsOptions,
) (*FundProvider, error) {
	v := validator.New()

	v.Check(!balance.IsZero(), "balance", "balance is required")
	v.Check(!fundProviderType.IsZero(), "fundProviderType", "fundProviderType is required")

	if err := v.Err(); err != nil {
		return nil, err
	}

	providerDetails, err := NewProviderDetails(fundProviderType, options)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &FundProvider{
		id:               id,
		balance:          balance,
		currency:         balance.Currency(),
		fundProviderType: fundProviderType,
		allocation:       allocation,
		providerDetails:  providerDetails,
	}, nil
}

func (f *FundProvider) ID() uuid.UUID                      { return f.id }
func (f *FundProvider) Balance() valueobject.Money         { return f.balance }
func (f *FundProvider) Currency() valueobject.Currency     { return f.currency }
func (f *FundProvider) FundProviderType() FundProviderType { return f.fundProviderType }
func (f *FundProvider) Allocation() Allocation             { return f.allocation }
func (f *FundProvider) AggVersion() int32                  { return f.aggVerion }

func (f *FundProvider) AllocateToWallet(
	walletID uuid.UUID,
	amount valueobject.Money,
) error {
	newAllocation, err := f.allocation.Allocate(walletID, amount)
	if err != nil {
		return err
	}

	if newAllocation.totalAllocated.GreaterThan(f.balance) {
		return ErrInsufficientBalance
	}

	f.allocation = newAllocation

	return nil
}

func (f *FundProvider) TopUp(
	walletID uuid.UUID,
	amount valueobject.Money,
) error {
	// Update balance
	newBalance, err := f.balance.Add(amount)
	if err != nil {
		return err
	}

	// Update allocation
	newAllocation, err := f.allocation.IncreaseAllocation(walletID, amount)
	if err != nil {
		return err
	}

	if newAllocation.totalAllocated.GreaterThan(f.balance) {
		return ErrInsufficientBalance
	}

	f.balance = newBalance
	f.allocation = newAllocation

	return nil
}

func (f *FundProvider) Withdraw(
	walletID uuid.UUID,
	amount valueobject.Money,
) error {
	// Update balance
	newBalance, err := f.balance.Subtract(amount)
	if err != nil {
		return err
	}

	if newBalance.IsNegative() {
		return ErrInsufficientBalance
	}

	// Update Allocation
	newAllocation, err := f.allocation.DecreaseAllocation(walletID, amount)
	if err != nil {
		return err
	}

	f.balance = newBalance
	f.allocation = newAllocation

	return nil
}
