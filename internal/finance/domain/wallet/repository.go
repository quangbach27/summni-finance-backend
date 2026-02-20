package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(
		ctx context.Context,
		wID uuid.UUID,
		fpIDs ...uuid.UUID,
	) (*Wallet, error)
	Create(ctx context.Context, wallet *Wallet) error
	Update(ctx context.Context, wallet *Wallet) error
}
