package valueobject_test

import (
	"sumni-finance-backend/internal/common/valueobject"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCurrency(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCode    string
		expectError bool
	}{
		{
			name:        "Success: USD",
			input:       "USD",
			wantCode:    "USD",
			expectError: false,
		},
		{
			name:        "Success: VND",
			input:       "VND",
			wantCode:    "VND",
			expectError: false,
		},
		{
			name:        "Success: Lowercase input normalized",
			input:       "usd",
			wantCode:    "USD",
			expectError: false,
		},
		{
			name:        "Success: Input with whitespace",
			input:       "  krw  ",
			wantCode:    "KRW",
			expectError: false,
		},
		{
			name:        "Failure: Unsupported currency",
			input:       "UNK",
			expectError: true,
		},
		{
			name:        "Failure: Empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueobject.NewCurrency(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, got.IsZero(), "Expected zero value on error")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCode, got.Code())
				assert.False(t, got.IsZero())
			}
		})
	}
}

func TestCurrency_Equal(t *testing.T) {
	tests := []struct {
		name  string
		base  valueobject.Currency
		other valueobject.Currency
		want  bool
	}{
		{
			name:  "True: Same currency instance",
			base:  valueobject.USD,
			other: valueobject.USD,
			want:  true,
		},
		{
			name:  "False: Different currencies",
			base:  valueobject.USD,
			other: valueobject.VND,
			want:  false,
		},
		{
			name:  "False: Compare with Zero value",
			base:  valueobject.USD,
			other: valueobject.Currency{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.base.Equal(tt.other)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCurrency_IsZero(t *testing.T) {
	t.Run("True for zero value", func(t *testing.T) {
		c := valueobject.Currency{}
		assert.True(t, c.IsZero())
	})

	t.Run("False for initialized value", func(t *testing.T) {
		c := valueobject.USD
		assert.False(t, c.IsZero())
	})
}
