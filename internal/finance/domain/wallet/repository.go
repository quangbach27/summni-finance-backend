package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(
		ctx context.Context,
		wID uuid.UUID,
	) (*Wallet, error)

	GetByIDWithProviders(
		ctx context.Context,
		wID uuid.UUID,
		spec ProviderAllocationSpec,
	) (*Wallet, error)

	Create(ctx context.Context, wallet *Wallet) error

	Update(
		ctx context.Context,
		wID uuid.UUID,
		spec ProviderAllocationSpec,
		updateFunc func(*Wallet) error,
	) error
}
