package command

import (
	"context"
	"fmt"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
)

type CreateWalletCmd struct {
	Currency    string
	Allocations []CreateWalletCmdAllocation
}
type CreateWalletCmdAllocation struct {
	ProviderID uuid.UUID
	Allocated  int64
}

type CreateWalletHandler cqrs.CommandHandler[CreateWalletCmd]

type createWalletHandler struct {
	walletRepo       wallet.Repository
	fundProviderRepo fundprovider.Repository
}

func NewCreateWalletHandler(
	walletRepo wallet.Repository,
	fundProviderRepo fundprovider.Repository,
) createWalletHandler {
	return createWalletHandler{
		walletRepo:       walletRepo,
		fundProviderRepo: fundProviderRepo,
	}
}

func (h createWalletHandler) Handle(ctx context.Context, cmd CreateWalletCmd) error {
	currency, err := valueobject.NewCurrency(cmd.Currency)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-currency")
	}

	wallet, err := wallet.NewWallet(currency)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-wallet")
	}

	for _, allocation := range cmd.Allocations {
		fundProvider, err := h.fundProviderRepo.GetByID(ctx, allocation.ProviderID)
		if err != nil {
			return httperr.NewUnknowError(err, "failed-to-retrieve-fund-provider")
		}

		if fundProvider == nil {
			return httperr.NewIncorrectInputError(
				fmt.Errorf("fund provider is not found with id: %s", allocation.ProviderID.String()),
				"fund-provider-not-found",
			)
		}

		allocated, err := valueobject.NewMoney(allocation.Allocated, currency)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "invalid-allocated")
		}

		if err = wallet.AddFundProvider(fundProvider, allocated); err != nil {
			return httperr.NewIncorrectInputError(err, "failed-to-allocated")
		}
	}

	return h.walletRepo.Create(ctx, wallet)
}
