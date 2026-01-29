package valueobject_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name        string
		amount      int64
		currency    valueobject.Currency
		expectError bool
	}{
		{
			name:        "Success: Valid positive amount",
			amount:      100,
			currency:    valueobject.USD,
			expectError: false,
		},
		{
			name:        "Success: Zero amount is allowed",
			amount:      0,
			currency:    valueobject.VND,
			expectError: false,
		},
		{
			name:        "Success: Negative amount",
			amount:      -100,
			currency:    valueobject.KRW,
			expectError: false,
		},
		{
			name:        "Failure: Missing currency",
			amount:      100,
			currency:    valueobject.Currency{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueobject.NewMoney(tt.amount, tt.currency)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.amount, got.Amount())
				assert.Equal(t, tt.currency, got.Currency())
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	// Setup fixtures
	usd100, _ := valueobject.NewMoney(100, valueobject.USD)
	usd50, _ := valueobject.NewMoney(50, valueobject.USD)
	usd0, _ := valueobject.NewMoney(0, valueobject.USD)
	krw1000, _ := valueobject.NewMoney(1000, valueobject.KRW)

	tests := []struct {
		name        string
		base        valueobject.Money
		other       valueobject.Money
		wantAmount  int64
		expectError bool
	}{
		{
			name:        "Success: Normal addition",
			base:        usd100,
			other:       usd50,
			wantAmount:  150,
			expectError: false,
		},
		{
			name:        "Success: Add zero",
			base:        usd100,
			other:       usd0,
			wantAmount:  100,
			expectError: false,
		},
		{
			name:        "Failure: Different currencies",
			base:        usd100,
			other:       krw1000,
			wantAmount:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.base.Add(tt.other)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if got.Amount() != tt.wantAmount {
					t.Errorf("Add() amount = %d, want %d", got.Amount(), tt.wantAmount)
				}
				assert.Equal(t, tt.wantAmount, got.Amount())
				assert.True(t, got.Currency().Equal(tt.base.Currency()))
			}
		})
	}
}

func TestMoney_Subtract(t *testing.T) {
	// Setup fixtures
	usd150, _ := valueobject.NewMoney(150, valueobject.USD)
	usd100, _ := valueobject.NewMoney(100, valueobject.USD)
	usd50, _ := valueobject.NewMoney(50, valueobject.USD)
	usd0, _ := valueobject.NewMoney(0, valueobject.USD)
	krw1000, _ := valueobject.NewMoney(1000, valueobject.KRW)

	tests := []struct {
		name        string
		base        valueobject.Money
		other       valueobject.Money
		wantAmount  int64
		expectError bool
	}{
		{
			name:        "Success: Normal subtraction",
			base:        usd100,
			other:       usd50,
			wantAmount:  50,
			expectError: false,
		},
		{
			name:        "Success: Subtract to zero",
			base:        usd100,
			other:       usd100,
			wantAmount:  0,
			expectError: false,
		},
		{
			name:        "Success: Subtract zero",
			base:        usd100,
			other:       usd0,
			wantAmount:  100,
			expectError: false,
		},
		{
			name:        "Sucess: Result negative",
			base:        usd100,
			other:       usd150,
			wantAmount:  -50,
			expectError: false,
		},
		{
			name:        "Failure: Different currencies",
			base:        usd100,
			other:       krw1000,
			wantAmount:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.base.Subtract(tt.other)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantAmount, got.Amount())
			}
		})
	}
}

func TestMoney_LessOrEqualThan(t *testing.T) {
	// Setup fixtures
	usd150, _ := valueobject.NewMoney(150, valueobject.USD)
	usd100, _ := valueobject.NewMoney(100, valueobject.USD)
	krw1000, _ := valueobject.NewMoney(1000, valueobject.KRW)

	tests := []struct {
		name  string
		base  valueobject.Money
		other valueobject.Money
		want  bool
	}{
		{
			name:  "True: 100 <= 150",
			base:  usd100,
			other: usd150,
			want:  true,
		},
		{
			name:  "True: 100 <= 100",
			base:  usd100,
			other: usd100,
			want:  true,
		},
		{
			name:  "False: 150 <= 100",
			base:  usd150,
			other: usd100,
			want:  false,
		},
		{
			name:  "False: Different currencies",
			base:  usd100,
			other: krw1000,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.LessOrEqualThan(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMoney_LessThan(t *testing.T) {
	// Setup fixtures
	usd150, _ := valueobject.NewMoney(150, valueobject.USD)
	usd100, _ := valueobject.NewMoney(100, valueobject.USD)
	krw1000, _ := valueobject.NewMoney(1000, valueobject.KRW)

	tests := []struct {
		name  string
		base  valueobject.Money
		other valueobject.Money
		want  bool
	}{
		{
			name:  "True: 100 < 150",
			base:  usd100,
			other: usd150,
			want:  true,
		},
		{
			name:  "False: 100 < 100",
			base:  usd100,
			other: usd100,
			want:  false,
		},
		{
			name:  "False: 150 < 100",
			base:  usd150,
			other: usd100,
			want:  false,
		},
		{
			name:  "False: Different currencies",
			base:  usd100,
			other: krw1000,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.LessThan(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	// Setup fixtures
	usd150, _ := valueobject.NewMoney(150, valueobject.USD)
	usd100, _ := valueobject.NewMoney(100, valueobject.USD)
	krw1000, _ := valueobject.NewMoney(1000, valueobject.KRW)

	tests := []struct {
		name  string
		base  valueobject.Money
		other valueobject.Money
		want  bool
	}{
		{
			name:  "True: 150 > 100",
			base:  usd150,
			other: usd100,
			want:  true,
		},
		{
			name:  "False: 100 > 100",
			base:  usd100,
			other: usd100,
			want:  false,
		},
		{
			name:  "False: 100 > 150",
			base:  usd100,
			other: usd150,
			want:  false,
		},
		{
			name:  "False: Different currencies",
			base:  usd100,
			other: krw1000,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.GreaterThan(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMoney_GreaterOrEqualThan(t *testing.T) {
	// Setup fixtures
	usd150, _ := valueobject.NewMoney(150, valueobject.USD)
	usd100, _ := valueobject.NewMoney(100, valueobject.USD)
	krw1000, _ := valueobject.NewMoney(1000, valueobject.KRW)

	tests := []struct {
		name  string
		base  valueobject.Money
		other valueobject.Money
		want  bool
	}{
		{
			name:  "True: 150 >= 100",
			base:  usd150,
			other: usd100,
			want:  true,
		},
		{
			name:  "True: 100 >= 100",
			base:  usd100,
			other: usd100,
			want:  true,
		},
		{
			name:  "False: 100 >= 150",
			base:  usd100,
			other: usd150,
			want:  false,
		},
		{
			name:  "False: Different currencies",
			base:  usd100,
			other: krw1000,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.GreaterOrEqualThan(tt.other)

			assert.Equal(t, tt.want, got)
		})
	}
}
