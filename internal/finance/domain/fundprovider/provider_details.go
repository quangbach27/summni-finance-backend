package fundprovider

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
)

type BankProviderDetails struct {
	bankName      string
	accountOwner  string
	accountNumber string
}

func NewBankProviderDetails(bankName, accountOwner, accountNumber string) (BankProviderDetails, error) {
	v := validator.New()

	v.Required(bankName, "bankName")
	v.Required(accountOwner, "accountOwner")
	v.Required(accountNumber, "accountNumber")

	if err := v.Err(); err != nil {
		return BankProviderDetails{}, err
	}

	return BankProviderDetails{
		bankName:      bankName,
		accountOwner:  accountOwner,
		accountNumber: accountNumber,
	}, nil
}

func (b BankProviderDetails) BankName() string { return b.bankName }

func (b BankProviderDetails) AccountOwner() string { return b.accountOwner }

func (b BankProviderDetails) AccountNumber() string { return b.accountNumber }

type CashProviderDetails struct {
	name string
}

func NewCashProviderDetails(name string) (CashProviderDetails, error) {
	if name == "" {
		return CashProviderDetails{}, errors.New("name is required")
	}

	return CashProviderDetails{
		name: name,
	}, nil
}

func (c CashProviderDetails) Name() string { return c.name }

type ProviderDetails struct {
	bankDetails BankProviderDetails
	cashDetails CashProviderDetails
}

type ProviderDetailsOptions struct {
	// Cash
	Name string

	// Bank
	BankName      string
	AccountOwner  string
	AccountNumber string
}

func NewProviderDetails(providerType FundProviderType, options ProviderDetailsOptions) (ProviderDetails, error) {
	if providerType.IsCash() {
		cashDetails, err := NewCashProviderDetails(options.Name)
		if err != nil {
			return ProviderDetails{}, err
		}

		return ProviderDetails{
			cashDetails: cashDetails,
		}, nil
	}

	if providerType.IsBank() {
		bankDetails, err := NewBankProviderDetails(
			options.BankName,
			options.AccountOwner,
			options.AccountNumber,
		)
		if err != nil {
			return ProviderDetails{}, err
		}

		return ProviderDetails{
			bankDetails: bankDetails,
		}, nil
	}

	return ProviderDetails{}, errors.New("unknow error")
}

func (p ProviderDetails) BankDetails() BankProviderDetails { return p.bankDetails }

func (p ProviderDetails) CashDetails() CashProviderDetails { return p.cashDetails }
