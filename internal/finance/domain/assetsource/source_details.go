package assetsource

import "sumni-finance-backend/internal/common/validator"

// BankDetails: Value Object (Grouping related information)
type BankDetails struct {
	bankName      string
	accountNumber string
}

func NewBankDetails(bankName, accountNumber string) (BankDetails, error) {
	validator := validator.New()

	validator.Required(bankName, "bankName")
	validator.Required(accountNumber, "accountNumber")

	if err := validator.Err(); err != nil {
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
