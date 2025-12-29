package command

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/assetsource"
	"sumni-finance-backend/internal/finance/domain/wallet"
)

type CreateWalletCmd struct {
	Name         string
	CurrencyCode string
	IsStrictMode bool
	Allocations  []AllocationItem
}

type AllocationItem struct {
	AssetSourceID string
	Amount        int64
}

type CreateWalletHandler cqrs.CommandHandler[CreateWalletCmd]

type createWalletHandler struct {
	walletRepo      wallet.Repository
	assetSourceRepo assetsource.Repository
}

func NewCreateWalletHandler(
	walletRepo wallet.Repository,
	assetSourceRepo assetsource.Repository,
) CreateWalletHandler {
	return cqrs.ApplyCommandDecorators(&createWalletHandler{
		walletRepo:      walletRepo,
		assetSourceRepo: assetSourceRepo,
	})
}

func (handler *createWalletHandler) Handle(ctx context.Context, cmd CreateWalletCmd) error {
	if len(cmd.Allocations) == 0 {
		return httperr.NewIncorrectInputError(
			errors.New("wallet must have at least one allocation"),
			"missing-allocation",
		)
	}

	allocationDomainList, err := handler.buildAllocationListDomain(ctx, cmd.Allocations)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "fail-to-build-allocation")
	}

	walletDomain, err := handler.buildWalletDomain(cmd, allocationDomainList)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "fail-to-build-wallet")
	}

	if err = handler.walletRepo.Create(ctx, walletDomain); err != nil {
		return httperr.NewUnknowError(err, "persist-wallet-failed")
	}

	return nil
}

func (handler *createWalletHandler) buildAllocationListDomain(ctx context.Context, allocationItems []AllocationItem) ([]*wallet.Allocation, error) {
	allocationDomainList := make([]*wallet.Allocation, 0, len(allocationItems))

	for _, item := range allocationItems {
		assetSourceID, err := assetsource.NewID(item.AssetSourceID)
		if err != nil {
			return nil, err
		}

		// TODO: check this can be applied spec Design Pattern
		assetSource, err := handler.assetSourceRepo.GetByID(ctx, assetSourceID)
		if err != nil {
			return nil, err
		}

		amount, err := valueobject.NewMoney(item.Amount, assetSource.Currency())
		if err != nil {
			return nil, err
		}

		allocationDomain, err := wallet.NewAllocation(assetSourceID, amount)
		if err != nil {
			return nil, err
		}

		allocationDomainList = append(allocationDomainList, allocationDomain)
	}

	return allocationDomainList, nil
}

func (handler *createWalletHandler) buildWalletDomain(
	cmd CreateWalletCmd,
	allocationDomainList []*wallet.Allocation,
) (*wallet.Wallet, error) {
	walletCurrency, err := valueobject.NewCurrency(cmd.CurrencyCode)
	if err != nil {
		return nil, err
	}

	walletDomain, err := wallet.NewWallet(cmd.Name, walletCurrency, cmd.IsStrictMode, allocationDomainList)
	if err != nil {
		return nil, err
	}

	return walletDomain, nil
}
