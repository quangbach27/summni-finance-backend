package db

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type walletRepo struct {
	queries *store.Queries
	pgPool  *pgxpool.Pool
}

func NewWalletRepo(
	queries *store.Queries,
	pgPool *pgxpool.Pool,
) (*walletRepo, error) {
	if queries == nil || pgPool == nil {
		return nil, errors.New("missing dependencies")
	}

	return &walletRepo{
		queries: queries,
		pgPool:  pgPool,
	}, nil
}

func (r *walletRepo) GetByID(
	ctx context.Context,
	wID uuid.UUID,
	fpIDs ...uuid.UUID,
) (*wallet.Wallet, error) {
	return nil, nil
}

func (r *walletRepo) Create(ctx context.Context, wallet *wallet.Wallet) error {
	return r.queries.CreateWallet(ctx, store.CreateWalletParams{
		ID:       wallet.ID(),
		Balance:  wallet.Balance().Amount(),
		Currency: wallet.Currency().Code(),
		Version:  wallet.Version(),
	})
}

func (r *walletRepo) Update(ctx context.Context, wallet *wallet.Wallet) error {
	return nil
}
