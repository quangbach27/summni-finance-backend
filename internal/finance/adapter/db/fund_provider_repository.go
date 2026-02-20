package db

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/fundprovider"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type fundProviderRepo struct {
	queries *store.Queries
	pgPool  *pgxpool.Pool
}

func NewFundProviderRepo(
	queries *store.Queries,
	pgPool *pgxpool.Pool,
) (*fundProviderRepo, error) {
	if queries == nil || pgPool == nil {
		return nil, errors.New("missing dependencies")
	}

	return &fundProviderRepo{
		queries: queries,
		pgPool:  pgPool,
	}, nil
}

func (r *fundProviderRepo) Create(
	ctx context.Context,
	fundProvider *fundprovider.FundProvider,
) error {
	return r.queries.CreateFundProvider(ctx, store.CreateFundProviderParams{
		ID:                fundProvider.ID(),
		Balance:           fundProvider.Balance().Amount(),
		Currency:          fundProvider.Currency().Code(),
		UnallocatedAmount: fundProvider.UnallocatedBalance().Amount(),
		Version:           fundProvider.Version(),
	})
}

func (r *fundProviderRepo) GetByID(ctx context.Context, fpID uuid.UUID) (*fundprovider.FundProvider, error) {
	return nil, nil
}

func (r *fundProviderRepo) GetByIDs(ctx context.Context, fpID []uuid.UUID) ([]*fundprovider.FundProvider, error) {
	return nil, nil
}
