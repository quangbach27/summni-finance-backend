package wallet

import (
	"errors"
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
	if assetSourceID == assetsource.ID(uuid.Nil) {
		return nil, errors.New("assertSouceID is required")
	}

	if amount.IsZero() {
		return nil, errors.New("amount is required")
	}

	return &Allocation{
		assetSourceID: assetSourceID,
		amount:        amount,
	}, nil
}

func (a *Allocation) AssetSourceID() assetsource.ID { return a.assetSourceID }
func (a *Allocation) Amount() valueobject.Money     { return a.amount }
