package ledger

import (
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

type LedgerAccount struct {
	id uuid.UUID

	walletID     uuid.UUID
	transactions []Transaction
}

func NewLedgerAccount() (LedgerAccount, error)

type Transaction struct {
	id              string
	walletBalance   valueobject.Money
	providerID      uuid.UUID
	providerBalance valueobject.Money
	description     string
	remitter        string
	remitterBank    string
}
