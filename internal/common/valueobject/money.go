package valueobject

import (
	"errors"
)

type Money struct {
	amount   int64
	currency Currency
}

// NewMoney creates a new Money instance.
// It enforces that money cannot be negative at creation.
func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("money amount cannot be negative")
	}
	if currency.IsZero() {
		return Money{}, errors.New("money currency is required")
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
func (m Money) IsZero() bool {
	return m == Money{}
}

// Add sums two money objects.
func (m Money) Add(other Money) (Money, error) {
	// 1. Allow adding Zero (No error needed)
	if other.IsZero() {
		return m, nil
	}

	// 2. Validate Currency
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add different currencies")
	}

	// 3. Optional: Check for Integer Overflow (Int64 is huge, but good for safety)
	// total := m.amount + other.amount
	// if total < m.amount { return Money{}, errors.New("integer overflow") }

	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract deducts money.
func (m Money) Subtract(other Money) (Money, error) {
	// 1. Allow subtracting Zero
	if other.IsZero() {
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
// Note: It returns false if currencies mismatch (or you could panic/error).
func (m Money) LessOrEqualThan(other Money) bool {
	if m.currency != other.currency {
		return false // Or handle mismatch differently
	}

	// Fixed Logic: Return true if m <= other
	return m.amount <= other.amount
}
