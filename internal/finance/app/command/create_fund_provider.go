package command

import (
	"context"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
)

type CreateFundProviderCmd struct {
	Balance  int64
	Currency string
}

type CreateFundProviderHandler cqrs.CommandHandler[CreateFundProviderCmd]

type createFundProviderHandler struct {
	fundProviderRepo fundprovider.Repository
}

func NewCreateFundProviderHandler(fundProviderRepo fundprovider.Repository) CreateFundProviderHandler {
	return createFundProviderHandler{
		fundProviderRepo: fundProviderRepo,
	}
}

func (h createFundProviderHandler) Handle(ctx context.Context, cmd CreateFundProviderCmd) error {
	currency, err := valueobject.NewCurrency(cmd.Currency)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-currency")
	}

	balance, err := valueobject.NewMoney(cmd.Balance, currency)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-balance")
	}

	fundProvider, err := fundprovider.NewFundProvider(balance)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-fund-provider")
	}

	err = h.fundProviderRepo.Create(ctx, fundProvider)
	if err != nil {
		return httperr.NewUnknowError(err, "failed-to-create-fund-provider")
	}

	return nil
}
