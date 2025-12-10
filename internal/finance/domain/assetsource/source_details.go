package assetsource

import (
	common_errors "sumni-finance-backend/internal/common/errors"
)

// BankDetails: Value Object (Grouping related information)
type BankDetails struct {
	bankName      string
	accountNumber string
}

func NewBankDetails(bankName, accountNumber string) (BankDetails, error) {
	validationErrs := &common_errors.ValidationErrors{}
	if bankName == "" {
		validationErrs.Add("bankName", "bank name is required")
	}

	if accountNumber == "" {
		validationErrs.Add("accountNumber", "account number is required")
	}

	if err := validationErrs.AsError(); err != nil {
		return BankDetails{}, err
	}

	return BankDetails{
		bankName:      bankName,
		accountNumber: accountNumber,
	}, nil
}

func (b BankDetails) BankName() string      { return b.bankName }
func (b BankDetails) AccountNumber() string { return b.accountNumber }
func (b BankDetails) IsZero() bool          { return b == BankDetails{} }
