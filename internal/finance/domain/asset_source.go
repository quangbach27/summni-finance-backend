package domain

import (
	"errors"
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

type AssetSourceID uuid.UUID

var (
	TypeBank = SourceType{code: "BANK"}
	TypeCash = SourceType{code: "CASH"}
)

type SourceType struct{ code string }

// BankDetails: Value Object (Grouping related information)
type BankDetails struct {
	bankName      string
	accountNumber string
}

// NewBankDetails creates the required details struct with validation.
func NewBankDetails(bankName, accountNumber string) (BankDetails, error) {
	if bankName == "" {
		return BankDetails{}, errors.New("bank name is required")
	}
	if accountNumber == "" {
		return BankDetails{}, errors.New("account number is required")
	}
	return BankDetails{
		bankName:      bankName,
		accountNumber: accountNumber,
	}, nil
}
func (b BankDetails) BankName() string      { return b.bankName }
func (b BankDetails) AccountNumber() string { return b.accountNumber }

// AssetSource: The main Aggregate Root
type AssetSource struct {
	id         AssetSourceID
	name       string
	sourceType SourceType // Struct (TypeBank or TypeCash)
	balance    valueobject.Money
	ownerID    uuid.UUID

	bankDetails *BankDetails // Nil if SourceType is Cash
}

// --- FACTORY 1: BANK ASSET ---

// NewBankAssetSource creates an AssetSource instance strictly typed as BANK.
// Enforces invariants: Bank accounts must have details.
func NewBankAssetSource(
	ownerID uuid.UUID,
	name string,
	initBalance valueobject.Money,
	bankName string,
	accountNumber string,
) (*AssetSource, error) {
	if name == "" {
		return nil, errors.New("asset name is required")
	}

	// 1. Validate and create BankDetails Value Object
	details, err := NewBankDetails(bankName, accountNumber)
	if err != nil {
		return nil, err
	}

	// 2. Generate ID and initial balance
	id, _ := uuid.NewV7()

	return &AssetSource{
		id:         AssetSourceID(id),
		ownerID:    ownerID,
		name:       name,
		sourceType: TypeBank, // Enforced type
		balance:    initBalance,
		// 3. Assign the details pointer
		bankDetails: &details,
	}, nil
}

// --- FACTORY 2: CASH ASSET ---

// NewCashAssetSource creates an AssetSource instance strictly typed as CASH.
// Invariant: Cash assets cannot have bank details.
func NewCashAssetSource(
	ownerID uuid.UUID,
	name string,
	initBalance valueobject.Money,
) (*AssetSource, error) {
	if name == "" {
		return nil, errors.New("asset name is required")
	}

	id, _ := uuid.NewV7()

	return &AssetSource{
		id:         AssetSourceID(id),
		ownerID:    ownerID,
		name:       name,
		sourceType: TypeCash, // Enforced type
		balance:    initBalance,

		// CRUCIAL: Must be nil for Cash to maintain integrity
		bankDetails: nil,
	}, nil
}

func (a *AssetSource) ID() AssetSourceID          { return a.id }
func (a *AssetSource) Name() string               { return a.name }
func (a *AssetSource) OwnerID() uuid.UUID         { return a.ownerID }
func (a *AssetSource) Type() SourceType           { return a.sourceType }
func (a *AssetSource) Balance() valueobject.Money { return a.balance }

// GetBankDetails provides read-only access to Bank specifics, returns nil if Cash Asset.
func (a *AssetSource) GetBankDetails() *BankDetails {
	return a.bankDetails
}
