package db

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/fundprovider"

	"github.com/google/uuid"
)

type fundProviderRepository struct {
	queries *store.Queries
}

func NewFundProviderRepository(queries *store.Queries) (*fundProviderRepository, error) {
	if queries == nil {
		return nil, errors.New("queries is required")
	}

	return &fundProviderRepository{
		queries: queries,
	}, nil
}

func (r *fundProviderRepository) GetByID(ctx context.Context, fpID uuid.UUID) (*fundprovider.FundProvider, error) {
	return nil, nil
}

func (r *fundProviderRepository) Create(ctx context.Context, fundProvider *fundprovider.FundProvider) error {
	params := store.CreateFundProviderParams{
		ID:       fundProvider.ID(),
		Balance:  fundProvider.Balance().Amount(),
		Currency: fundProvider.Balance().Currency().Code(),
		Version:  fundProvider.Verions(),
	}

	return r.queries.CreateFundProvider(ctx, params)
}
