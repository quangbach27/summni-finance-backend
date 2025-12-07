package valueobject

import (
	"errors"
	"strings"
)

type Currency struct {
	code string
}

var (
	USD = Currency{code: "USD"}
	VND = Currency{code: "VND"}
	KRW = Currency{code: "KRW"}

	supportedCurrencies = []Currency{USD, VND, KRW}
)

func NewCurrency(code string) (Currency, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	for _, supportedCurrency := range supportedCurrencies {
		if supportedCurrency.code == code {
			return supportedCurrency, nil
		}
	}

	return Currency{}, errors.New("unsupported currency: " + code)
}

func (c Currency) Code() string { return c.code }
func (c Currency) IsZero() bool { return c == Currency{} }

func (c Currency) Equal(other Currency) bool {
	return c.code == other.code
}
