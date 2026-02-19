package command

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
)

type AllocateFundProviderCmd struct {
	WalletID      uuid.UUID
	FundProviders []FundProviderCmd
}

type FundProviderCmd struct {
	ID        uuid.UUID
	Allocated int64
}

type AllocateFundProviderHandler cqrs.CommandHandler[AllocateFundProviderCmd]

type allocateFundProviderHandler struct {
	fundProviderRepo fundprovider.Repository
	walletRepo       wallet.Repository
}

func NewAllocateFundProviderHandler(fundProviderRepo fundprovider.Repository, walletRepo wallet.Repository) *allocateFundProviderHandler {
	return &allocateFundProviderHandler{
		fundProviderRepo: fundProviderRepo,
		walletRepo:       walletRepo,
	}
}

func (h *allocateFundProviderHandler) Handle(ctx context.Context, cmd AllocateFundProviderCmd) error {
	err := h.validateComand(cmd)
	if err != nil {
		return err
	}

	walletDomain, err := h.walletRepo.GetByID(ctx, cmd.WalletID)
	if err != nil {
		return httperr.NewUnknowError(err, "failed-to-retrieve-wallet")
	}
	if walletDomain == nil {
		return httperr.NewIncorrectInputError(fmt.Errorf("wallet does not exist id: %s", cmd.WalletID.String()), "wallet-does-not-exist")
	}

	fpIDs := make([]uuid.UUID, 0, len(cmd.FundProviders))
	for _, fp := range cmd.FundProviders {
		fpIDs = append(fpIDs, fp.ID)
	}

	fundProvidersDomain, err := h.fundProviderRepo.GetByIDs(ctx, fpIDs)
	if err != nil {
		return httperr.NewUnknowError(err, "failed-to-retrieve-fund-providers")
	}
	if len(fpIDs) != len(fundProvidersDomain) {
		return httperr.NewIncorrectInputError(errors.New("Some provider are missing"), "fund-provider-missing")
	}

	for _, fundProvider := range fundProvidersDomain {
		allocatedMoney, err := h.getAllocatedAmountByFundProviderID(
			fundProvider.ID(),
			cmd.FundProviders,
			walletDomain.Balance().Currency(),
		)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "invalid_allocated_money")
		}

		err = walletDomain.AddFundProvider(fundProvider, allocatedMoney)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "failed_to_allocated_fund_provider")
		}
	}

	if err = h.walletRepo.Update(ctx, walletDomain); err != nil {
		return httperr.NewUnknowError(err, "failed-to-update-wallet")
	}

	return nil
}

func (h *allocateFundProviderHandler) validateComand(cmd AllocateFundProviderCmd) error {
	if cmd.WalletID == uuid.Nil {
		return httperr.NewIncorrectInputError(errors.New("Wallet id is required"), "missing-wallet-id")
	}

	if len(cmd.FundProviders) == 0 {
		return httperr.NewIncorrectInputError(errors.New("FundProviders must not empty"), "missing-fund-providers")
	}

	return nil
}

func (h *allocateFundProviderHandler) getAllocatedAmountByFundProviderID(
	fpID uuid.UUID,
	fundProvidersCmd []FundProviderCmd,
	currency valueobject.Currency,
) (valueobject.Money, error) {
	var allocated int64

	for _, fundProviderCmd := range fundProvidersCmd {
		if fundProviderCmd.ID == fpID {
			allocated = fundProviderCmd.Allocated
		}
	}

	money, err := valueobject.NewMoney(allocated, currency)
	if err != nil {
		return valueobject.Money{}, err
	}

	return money, nil
}
