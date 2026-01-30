package fundprovider

import (
	"fmt"
	"strings"
)

var (
	CashProviderType = FundProviderType{value: "CASH"}
	BankProviderType = FundProviderType{value: "BANK"}
)

type ErrUnsupportedFundProviderType struct {
	ProviderType string
}

func (err ErrUnsupportedFundProviderType) Error() string {
	return fmt.Sprintf("'%s' is not supported", err.ProviderType)
}

type FundProviderType struct {
	value string
}

func FundProviderTypeOf(providerType string) (FundProviderType, error) {
	cleanStr := strings.ToUpper(strings.TrimSpace(providerType))

	if cleanStr == CashProviderType.value {
		return CashProviderType, nil
	}

	if cleanStr == BankProviderType.value {
		return BankProviderType, nil
	}

	return FundProviderType{}, ErrUnsupportedFundProviderType{ProviderType: cleanStr}
}

func (t FundProviderType) Value() string { return t.value }
func (t FundProviderType) IsZero() bool  { return t.value == "" }
func (t FundProviderType) IsCash() bool  { return t.value == CashProviderType.value }
func (t FundProviderType) IsBank() bool  { return t.value == BankProviderType.value }
