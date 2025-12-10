package assetsource

import (
	"errors"
	"fmt"
	"strings"
	commons_errors "sumni-finance-backend/internal/common/errors"
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

type ID uuid.UUID

// AssetSource: Aggregate Root
type AssetSource struct {
	id         ID
	balance    valueobject.Money
	ownerID    uuid.UUID
	sourceType SourceType

	bankDetails *BankDetails // Nil if SourceType is Cash
}

func (as *AssetSource) ID() ID                     { return as.id }
func (as *AssetSource) Balance() valueobject.Money { return as.balance }
func (as *AssetSource) OwnerID() uuid.UUID         { return as.ownerID }
func (as *AssetSource) Type() SourceType           { return as.sourceType }
func (as *AssetSource) BankDetails() *BankDetails  { return as.bankDetails }

// newBaseAssetSource: Private template for shared asset initialization logic
func newBaseAssetSource(
	ownerID uuid.UUID,
	amount int64,
	currency valueobject.Currency,
	sourceType SourceType,
) (*AssetSource, error) {
	validationErrs := &commons_errors.ValidationErrors{}

	if ownerID == uuid.Nil {
		validationErrs.Add("ownerID", "ownerID is required")
	}

	if currency.IsZero() {
		validationErrs.Add("currency", "currency is required")
	}

	if sourceType.IsZero() {
		validationErrs.Add("sourceType", "sourceType is required")
	}

	initbalance, err := valueobject.NewMoney(amount, currency)
	if err != nil {
		validationErrs.Add("balance", err.Error())
	}

	if err := validationErrs.AsError(); err != nil {
		return nil, err
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate AssetSourceID: %w", err)
	}

	return &AssetSource{
		id:         ID(newID),
		balance:    initbalance,
		ownerID:    ownerID,
		sourceType: sourceType,
	}, nil
}

// --- Factory 1: BANK ASSET ---
func NewBankAssetSource(
	ownerID uuid.UUID,
	initAmount int64,
	currency valueobject.Currency,
	bankName string,
	accountNumber string,
) (*AssetSource, error) {
	validationErrs := &commons_errors.ValidationErrors{}

	assetSource, err := newBaseAssetSource(ownerID, initAmount, currency, BankType)

	if err != nil {
		var baseValidationErrs *commons_errors.ValidationErrors
		if errors.As(err, &baseValidationErrs) {
			validationErrs.Merge(baseValidationErrs)
		} else {
			return nil, err
		}
	}

	bankDetails, err := NewBankDetails(bankName, accountNumber)
	if err != nil {
		validationErrs.Add("bankDetails", err.Error())
	}
	if err = validationErrs.AsError(); err != nil {
		return nil, err
	}

	assetSource.bankDetails = &bankDetails

	return assetSource, nil
}

// --- Factory 2: CASH ASSET ---
func NewCashAssetSource(
	ownerID uuid.UUID,
	initAmount int64,
	currency valueobject.Currency,
) (*AssetSource, error) {
	assetSource, err := newBaseAssetSource(ownerID, initAmount, currency, CashType)
	if err != nil {
		return nil, err
	}

	return assetSource, nil
}

// --- Source Type (Sealed Enum) ---
type SourceType struct{ code string }

var (
	CashType = SourceType{code: "CASH"}
	BankType = SourceType{code: "BANK"}
)

func NewSourceTypeFromStr(soureTypeStr string) (SourceType, error) {
	code := strings.TrimSpace(strings.ToUpper(soureTypeStr))
	if code == CashType.code {
		return CashType, nil
	}

	if code == BankType.code {
		return BankType, nil
	}

	return SourceType{}, errors.New("unknow asset source type: " + soureTypeStr)
}

func (st SourceType) Code() string { return st.code }
func (st SourceType) IsZero() bool { return st == SourceType{} }
