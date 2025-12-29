package assetsource

import (
	"errors"
	"fmt"
	"strings"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID(idStr string) (ID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return ID{}, fmt.Errorf("invalid id format: %w", err)
	}

	return ID(id), nil
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

// AssetSource: Aggregate Root
type AssetSource struct {
	id         ID
	name       string
	balance    valueobject.Money
	ownerID    uuid.UUID
	sourceType SourceType
	officeID   uuid.UUID
	currency   valueobject.Currency

	bankDetails *BankDetails // Nil if SourceType is Cash
}

func (as *AssetSource) ID() ID                         { return as.id }
func (as *AssetSource) Name() string                   { return as.name }
func (as *AssetSource) Balance() valueobject.Money     { return as.balance }
func (as *AssetSource) OwnerID() uuid.UUID             { return as.ownerID }
func (as *AssetSource) Type() SourceType               { return as.sourceType }
func (as *AssetSource) BankDetails() *BankDetails      { return as.bankDetails }
func (as *AssetSource) OfficeID() uuid.UUID            { return as.officeID }
func (as *AssetSource) Currency() valueobject.Currency { return as.currency }

// newBaseAssetSource: Private template for shared asset initialization logic
func newBaseAssetSource(
	ownerID uuid.UUID,
	name string,
	amount int64,
	currency valueobject.Currency,
	sourceType SourceType,
	officeID uuid.UUID,
) (*AssetSource, error) {
	validator := validator.New()

	validator.Check(ownerID != uuid.Nil, "ownerID", "ownerID is required")
	validator.Required(name, "name")
	validator.Check(amount >= 0, "amount", "amount must be positive")
	validator.Check(!currency.IsZero(), "currency", "currency is required")
	validator.Check(!sourceType.IsZero(), "sourceType", "sourceType is required")
	validator.Check(officeID != uuid.Nil, "officeID", "ownerID is required")

	if err := validator.Err(); err != nil {
		return nil, err
	}

	initbalance, err := valueobject.NewMoney(amount, currency)
	if err != nil {
		return nil, err
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate AssetSourceID: %w", err)
	}

	return &AssetSource{
		id:         ID(newID),
		name:       name,
		balance:    initbalance,
		ownerID:    ownerID,
		sourceType: sourceType,
		officeID:   officeID,
		currency:   currency,
	}, nil
}

// --- Factory 1: BANK ASSET ---
func NewBankAssetSource(
	ownerID uuid.UUID,
	name string,
	initAmount int64,
	currency valueobject.Currency,
	bankName string,
	accountNumber string,
	officeID uuid.UUID,
) (*AssetSource, error) {
	validator := validator.New()

	assetSource, err := newBaseAssetSource(ownerID, name, initAmount, currency, BankType, officeID)
	if err != nil {
		if !validator.TryMerge(err) {
			return nil, err
		}
	}

	bankDetails, err := NewBankDetails(bankName, accountNumber)
	if err != nil {
		if !validator.TryMerge(err) {
			return nil, err
		}
	}

	if err = validator.Err(); err != nil {
		return nil, err
	}

	assetSource.bankDetails = &bankDetails

	return assetSource, nil
}

// --- Factory 2: CASH ASSET ---
func NewCashAssetSource(
	ownerID uuid.UUID,
	name string,
	initAmount int64,
	currency valueobject.Currency,
	officeID uuid.UUID,
) (*AssetSource, error) {
	assetSource, err := newBaseAssetSource(ownerID, name, initAmount, currency, CashType, officeID)
	if err != nil {
		return nil, err
	}

	return assetSource, nil
}

func UnmarshallFromDatabase(
	id uuid.UUID,
	ownerID uuid.UUID,
	name string,
	balance int64,
	sourceTypeStr string,
	currencyCode string,
	bankName string,
	accountNumber string,
	officeID uuid.UUID,
) (*AssetSource, error) {
	currency, err := valueobject.NewCurrency(currencyCode)
	if err != nil {
		return nil, err
	}

	balanceDomain, err := valueobject.NewMoney(balance, currency)
	if err != nil {
		return nil, err
	}

	sourceType, err := NewSourceTypeFromStr(sourceTypeStr)
	if err != nil {
		return nil, err
	}

	assetSource := AssetSource{
		id:         ID(id),
		name:       name,
		balance:    balanceDomain,
		ownerID:    ownerID,
		sourceType: sourceType,
		officeID:   officeID,
		currency:   currency,
	}

	if sourceType.IsBank() {
		bankDetails, err := NewBankDetails(bankName, accountNumber)
		if err != nil {
			return nil, err
		}

		assetSource.bankDetails = &bankDetails
	}

	return &assetSource, nil
}

// --- Source Type (Sealed Enum) ---
type SourceType struct{ code string }

var (
	CashType = SourceType{code: "CASH"}
	BankType = SourceType{code: "BANK"}
)

func NewSourceTypeFromStr(sourceTypeStr string) (SourceType, error) {
	code := strings.TrimSpace(strings.ToUpper(sourceTypeStr))
	if code == CashType.code {
		return CashType, nil
	}

	if code == BankType.code {
		return BankType, nil
	}

	return SourceType{}, errors.New("unknown asset source type: " + sourceTypeStr)
}

func (st SourceType) Code() string { return st.code }
func (st SourceType) IsZero() bool { return st == SourceType{} }
func (st SourceType) IsCash() bool { return st == CashType }
func (st SourceType) IsBank() bool { return st == BankType }
