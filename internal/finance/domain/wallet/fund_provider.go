package wallet

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

var (
	ErrCurrencyMismatch      = errors.New("currency mismatch")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrInsufficientAvailable = errors.New("insufficient available amount")
)

type FundProvider struct {
	id                           uuid.UUID
	balance                      valueobject.Money
	availableAmountForAllocation valueobject.Money
}

func NewFundProvider(
	balance valueobject.Money,
) (*FundProvider, error) {
	if balance.IsZero() {
		return nil, errors.New("balance is requried")
	}

	if balance.Amount() < 0 {
		return nil, errors.New("balance must be positive or equal zero")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &FundProvider{
		id:                           id,
		balance:                      balance,
		availableAmountForAllocation: balance,
	}, nil
}

func UnmarshallFundProviderFromDatabase(
	id uuid.UUID,
	balance valueobject.Money,
	availableAmountForAllocation valueobject.Money,
) (*FundProvider, error) {
	v := validator.New()

	v.Check(id != uuid.Nil, "id", "id is required")
	v.Check(!balance.IsZero(), "balance", "balance is required")
	v.Check(!availableAmountForAllocation.IsZero(), "availableAmountForAllocation", "availableAmountForAllocation is required")

	if err := v.Err(); err != nil {
		return nil, err
	}

	return &FundProvider{
		id:                           id,
		balance:                      balance,
		availableAmountForAllocation: availableAmountForAllocation,
	}, nil
}

func (p *FundProvider) ID() uuid.UUID              { return p.id }
func (p *FundProvider) Balance() valueobject.Money { return p.balance }
func (p *FundProvider) AvailableAmountForAllocation() valueobject.Money {
	return p.availableAmountForAllocation
}

func (p *FundProvider) TopUp(amount valueobject.Money) error {
	if !p.isAmountCurrencyValid(amount.Currency()) {
		return ErrCurrencyMismatch
	}

	newBalance, err := p.balance.Add(amount)
	if err != nil {
		return err
	}

	p.balance = newBalance
	return nil
}

func (p *FundProvider) Withdraw(amount valueobject.Money) error {
	if !p.isAmountCurrencyValid(amount.Currency()) {
		return ErrCurrencyMismatch
	}

	if amount.GreaterThan(p.balance) {
		return ErrInsufficientBalance
	}

	newBalance, err := p.balance.Subtract(amount)
	if err != nil {
		return err
	}

	p.balance = newBalance
	return nil
}

func (p *FundProvider) Allocate(
	allocatedAmount valueobject.Money,
) error {
	if allocatedAmount.GreaterThan(p.availableAmountForAllocation) {
		return ErrInsufficientAvailable
	}

	newAvailableAmount, err := p.availableAmountForAllocation.Subtract(allocatedAmount)
	if err != nil {
		return err
	}

	p.availableAmountForAllocation = newAvailableAmount
	return nil
}

func (p *FundProvider) isAmountCurrencyValid(currency valueobject.Currency) bool {
	return p.balance.Currency().Equal(currency)
}
