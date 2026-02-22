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
	if other.IsZero() {
		return Money{}, errors.New("empty input money when adding")
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

	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

// LessOrEqualThan compares amounts.
func (m Money) LessOrEqualThan(other Money) bool {
	if m.currency != other.currency {
		return false
	}

	return m.amount <= other.amount
}

// LessThan compares amounts.
func (m Money) LessThan(other Money) bool {
	if m.currency != other.currency {
		return false
	}

	return m.amount < other.amount
}

// GreaterThan compares amounts.
func (m Money) GreaterThan(other Money) bool {
	if m.currency != other.currency {
		return false
	}

	return m.amount > other.amount
}

// GreaterOrEqualThan compares amounts.
func (m Money) GreaterOrEqualThan(other Money) bool {
	if m.currency != other.currency {
		return false
	}

	return m.amount >= other.amount
}

func (m Money) Equal(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}
