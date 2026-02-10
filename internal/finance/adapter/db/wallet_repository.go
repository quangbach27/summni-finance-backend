package db

import (
	"context"
	"errors"
	commondb "sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type walletRepository struct {
	queries *store.Queries
	pgPool  *pgxpool.Pool
}

func NewWalletRepository(queries *store.Queries, pgPool *pgxpool.Pool) (*walletRepository, error) {
	if queries == nil {
		return nil, errors.New("queries is required")
	}

	if pgPool == nil {
		return nil, errors.New("pgPool is required")
	}

	return &walletRepository{
		queries: queries,
		pgPool:  pgPool,
	}, nil
}

func (r *walletRepository) GetByID(ctx context.Context, wID uuid.UUID) (*wallet.Wallet, error) {
	return nil, nil
}

func (r *walletRepository) Create(ctx context.Context, wallet *wallet.Wallet) error {
	// Start transaction
	tx, err := r.pgPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = commondb.FinishTransaction(ctx, tx, err)
	}()

	// Use queries with transaction
	qTx := r.queries.WithTx(tx)

	// 1. Create Wallet
	err = qTx.CreateWallet(ctx, store.CreateWalletParams{
		ID:       wallet.ID(),
		Balance:  wallet.Balance().Amount(),
		Currency: wallet.Balance().Currency().Code(),
		Version:  wallet.Version(),
	})
	if err != nil {
		return err
	}

	// 2. Batch create provider allocations
	providerAllocations := wallet.ProviderManager().GetFundProviderAllocations()
	if len(providerAllocations) == 0 {
		return nil
	}

	params := make([]store.CreatFundProviderAllocationParams, 0, len(providerAllocations))

	for _, allocation := range providerAllocations {
		params = append(params, store.CreatFundProviderAllocationParams{
			WalletID:        wallet.ID(),
			FundProviderID:  allocation.FundProvider().ID(),
			AllocatedAmount: allocation.Allocated().Amount(),
		})
	}

	rows, err := qTx.CreatFundProviderAllocation(ctx, params)
	if err != nil {
		return err
	}

	if rows != int64(len(params)) {
		return errors.New("copyfrom: partial fund_provider_allocation insert")
	}

	return nil
}
