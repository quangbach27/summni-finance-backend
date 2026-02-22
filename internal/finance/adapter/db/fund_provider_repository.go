package db

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/fundprovider"

	"github.com/google/uuid"
)

type fundProviderRepo struct {
	queries *store.Queries
}

func NewFundProviderRepo(
	queries *store.Queries,
) (*fundProviderRepo, error) {
	if queries == nil {
		return nil, errors.New("missing dependencies")
	}

	return &fundProviderRepo{
		queries: queries,
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
	fundProviderModel, err := r.queries.GetFundProviderByID(ctx, fpID)
	if err != nil {
		return nil, err
	}

	return fundprovider.UnmarshalFundProviderFromDatabase(
		fundProviderModel.ID,
		fundProviderModel.Balance,
		fundProviderModel.UnallocatedAmount,
		fundProviderModel.Currency,
		fundProviderModel.Version,
	)
}

func (r *fundProviderRepo) GetByIDs(ctx context.Context, fpID []uuid.UUID) ([]*fundprovider.FundProvider, error) {
	fpModels, err := r.queries.GetFundProvidersByIDs(ctx, fpID)
	if err != nil {
		return nil, err
	}

	fps := make([]*fundprovider.FundProvider, 0, len(fpModels))
	for _, model := range fpModels {
		fp, err := fundprovider.UnmarshalFundProviderFromDatabase(
			model.ID,
			model.Balance,
			model.UnallocatedAmount,
			model.Currency,
			model.Version,
		)
		if err != nil {
			return nil, err
		}
		fps = append(fps, fp)
	}

	return fps, nil
}
