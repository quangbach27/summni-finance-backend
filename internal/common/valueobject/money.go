package valueobject

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
)

type Money struct {
	amount   int64
	currency Currency
}

// NewMoney creates a new Money instance.
// It enforces that money cannot be negative at creation.
func NewMoney(amount int64, currency Currency) (Money, error) {
	validator := validator.New()

	validator.Check(amount >= 0, "money.amount", "amount cannot be negative")
	validator.Check(!currency.IsZero(), "money.currency", "currency is required")

	if err := validator.Err(); err != nil {
		return Money{}, err
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// Getters
func (m Money) Amount() int64      { return m.amount }
func (m Money) Currency() Currency { return m.currency }

// IsZero checks if the struct is the zero value (uninitialized)
func (m Money) IsZero() bool { return m == Money{} }

// Add sums two money objects.
func (m Money) Add(other Money) (Money, error) {
	// 1. Allow adding Zero (No error needed)
	if other.IsZero() || other.amount == 0 {
		return m, nil
	}

	// 2. Validate Currency
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add different currencies")
	}

	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract deducts money.
func (m Money) Subtract(other Money) (Money, error) {
	// 1. Allow subtracting Zero
	if other.IsZero() || other.amount == 0 {
		return m, nil
	}

	// 2. Validate Currency
	if m.currency != other.currency {
		return Money{}, errors.New("cannot subtract different currencies")
	}

	resultAmount := m.amount - other.amount

	// 3. Maintain Invariant: If NewMoney forbids negative, Subtract must also forbid it.
	if resultAmount < 0 {
		return Money{}, errors.New("insufficient funds: result cannot be negative")
	}

	return Money{
		amount:   resultAmount,
		currency: m.currency,
	}, nil
}

// LessOrEqualThan compares amounts.
func (m Money) LessOrEqualThan(other Money) bool {
	// Choosing to return false here implies they are not comparable.
	// Be careful not to rely on this for sorting mixed currencies.
	if m.currency != other.currency {
		return false
	}

	return m.amount <= other.amount
}
