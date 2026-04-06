package wallet

import (
	"context"
	"sumni-finance-backend/internal/finance/domain/ledger"

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

	GetByIDWithAccountingPeriod(
		ctx context.Context,
		wID uuid.UUID,
		yearMonth ledger.YearMonth,
	) (*Wallet, error)

	Create(ctx context.Context, wallet *Wallet) error

	CreateAllocations(
		ctx context.Context,
		wID uuid.UUID,
		allocationSpec ProviderAllocationSpec,
		allocatedFunc func(*Wallet) error,
	) error

	CreateTransactionRecords(
		ctx context.Context,
		wID uuid.UUID,
		allocationSpec ProviderAllocationSpec,
		yearMonth ledger.YearMonth,
		updateFunc func(w *Wallet) error,
	) error
}
