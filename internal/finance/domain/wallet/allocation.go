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
}

func NewAllocation(assetSourceID assetsource.ID, amount valueobject.Money) (*Allocation, error) {
	validator := validator.New()

	validator.Check(assetSourceID != assetsource.ID(uuid.Nil), "assetSourceID", "assetSourceID is required")
	validator.Check(!amount.IsZero(), "amount", "amount is required")

	if err := validator.Err(); err != nil {
		return nil, err
	}

	return &Allocation{
		assetSourceID: assetSourceID,
		amount:        amount,
	}, nil
}

func (a *Allocation) AssetSourceID() assetsource.ID { return a.assetSourceID }
func (a *Allocation) Amount() valueobject.Money     { return a.amount }
