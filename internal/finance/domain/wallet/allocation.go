package wallet

import (
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/assetsource"

	"github.com/google/uuid"
)

// --- ALLOCATION ---
type Allocation struct {
	assetSourceID assetsource.ID
	amount        valueobject.Money
	officeID      uuid.UUID
}

func NewAllocation(assetSourceID assetsource.ID, amount valueobject.Money, officeID uuid.UUID) (*Allocation, error) {
	validator := validator.New()

	validator.Check(assetSourceID != assetsource.ID(uuid.Nil), "assetSourceID", "assetSourceID is required")
	validator.Check(!amount.IsZero(), "amount", "amount is required")
	validator.Check(officeID != uuid.Nil, "officeID", "officeID is required")

	if err := validator.Err(); err != nil {
		return nil, err
	}

	return &Allocation{
		assetSourceID: assetSourceID,
		amount:        amount,
		officeID:      officeID,
	}, nil
}

func (a *Allocation) AssetSourceID() assetsource.ID { return a.assetSourceID }
func (a *Allocation) Amount() valueobject.Money     { return a.amount }
func (a *Allocation) OfficeID() uuid.UUID           { return a.officeID }
