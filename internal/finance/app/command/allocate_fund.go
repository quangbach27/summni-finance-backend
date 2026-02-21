package command

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
)

type AllocateFundCmd struct {
	WalletID  uuid.UUID
	Providers []AllocatedProviders
}

type AllocatedProviders struct {
	ID              uuid.UUID
	AllocatedAmount int64
}

type AllocateFundHandler cqrs.CommandHandler[AllocateFundCmd]

type allocateFundHandler struct {
	walletRepo       wallet.Repository
	fundProviderRepo fundprovider.Repository
}

func NewAllocateFundHandler(walletRepo wallet.Repository, fundProviderRepo fundprovider.Repository) *allocateFundHandler {
	return &allocateFundHandler{
		walletRepo:       walletRepo,
		fundProviderRepo: fundProviderRepo,
	}
}

func (h *allocateFundHandler) Handle(ctx context.Context, cmd AllocateFundCmd) error {
	if len(cmd.Providers) == 0 {
		return httperr.NewIncorrectInputError(
			errors.New("missing providers in command for allocation"),
			"invalid-providers",
		)
	}

	providerIDs := make([]uuid.UUID, 0, len(cmd.Providers))
	for _, p := range cmd.Providers {
		providerIDs = append(providerIDs, p.ID)
	}

	providerDomains, err := h.fundProviderRepo.GetByIDs(ctx, providerIDs)
	if err != nil {
		return httperr.NewUnknowError(err, "failed-to-retrieve-fund-provider")
	}

	err = h.walletRepo.Update(
		ctx,
		cmd.WalletID,
		wallet.NewAllocationBelongsToAnyProviderSpec(providerIDs),
		func(w *wallet.Wallet) error {
			for _, provider := range cmd.Providers {
				providerDomain := h.findFundProvider(providerDomains, provider.ID)
				if providerDomain == nil {
					return fmt.Errorf("fund provider '%s' did not existe", provider.ID)
				}

				err = w.AllocateFromFundProvider(providerDomain, provider.AllocatedAmount)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
	if err != nil {
		return httperr.NewUnknowError(err, "failed-to-allocate-fund")
	}

	return nil
}

func (h *allocateFundHandler) findFundProvider(
	fundProviders []*fundprovider.FundProvider,
	fpID uuid.UUID,
) *fundprovider.FundProvider {
	for _, provider := range fundProviders {
		if provider.ID() == fpID {
			return provider
		}
	}

	return nil
}
