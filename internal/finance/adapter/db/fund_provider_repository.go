package db

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/common/valueobject"
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
	model, err := r.queries.GetFundProviderByID(ctx, fpID)
	if err != nil {
		return nil, err
	}

	currency, err := valueobject.NewCurrency(model.Currency)
	if err != nil {
		return nil, err
	}

	balance, err := valueobject.NewMoney(model.Balance, currency)
	if err != nil {
		return nil, err
	}

	availableAmount, err := valueobject.NewMoney(int64(model.AvailableAmount), currency)
	if err != nil {
		return nil, err
	}

	fundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
		model.ID,
		balance,
		availableAmount,
		model.Version,
	)

	if err != nil {
		return nil, err
	}

	return fundProvider, nil
}

func (r *fundProviderRepository) Create(ctx context.Context, fundProvider *fundprovider.FundProvider) error {
	params := store.CreateFundProviderParams{
		ID:              fundProvider.ID(),
		Balance:         fundProvider.Balance().Amount(),
		Currency:        fundProvider.Balance().Currency().Code(),
		AvailableAmount: fundProvider.AvailableAmountForAllocation().Amount(),
		Version:         fundProvider.Verions(),
	}

	return r.queries.CreateFundProvider(ctx, params)
}

func (r *fundProviderRepository) GetByIDs(ctx context.Context, fpIDs []uuid.UUID) ([]*fundprovider.FundProvider, error) {
	return nil, nil
}
