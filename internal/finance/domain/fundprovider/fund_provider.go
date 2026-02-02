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

	providerDetails ProviderDetails
	version         int32
}

func NewFundProvider(
	balance valueobject.Money,
	fundProviderType FundProviderType,
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
		providerDetails:  providerDetails,
	}, nil
}

func (f *FundProvider) ID() uuid.UUID                      { return f.id }
func (f *FundProvider) Balance() valueobject.Money         { return f.balance }
func (f *FundProvider) Currency() valueobject.Currency     { return f.currency }
func (f *FundProvider) FundProviderType() FundProviderType { return f.fundProviderType }
func (f *FundProvider) Version() int32                     { return f.version }

func (f *FundProvider) TopUp(
	amount valueobject.Money,
) error {
	newBalance, err := f.balance.Add(amount)
	if err != nil {
		return err
	}

	f.balance = newBalance

	return nil
}

func (f *FundProvider) Withdraw(
	amount valueobject.Money,
) error {
	newBalance, err := f.balance.Subtract(amount)
	if err != nil {
		return err
	}

	if newBalance.IsNegative() {
		return ErrInsufficientBalance
	}

	f.balance = newBalance

	return nil
}
