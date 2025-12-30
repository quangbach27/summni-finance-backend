package command

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/assetsource"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
)

type CreateWalletResult struct {
	WalletID string `json:"walletId"`
}

type CreateWalletCmd struct {
	Name         string
	CurrencyCode string
	IsStrictMode bool
	OfficeID     string
	Allocations  []CreateWalletAllocation
}

type CreateWalletAllocation struct {
	AssetSourceID string
	Amount        int64
	OfficeID      string
}

type CreateWalletHandler cqrs.CommandHandler[CreateWalletCmd, CreateWalletResult]

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

func (handler *createWalletHandler) Handle(ctx context.Context, cmd CreateWalletCmd) (CreateWalletResult, error) {
	if len(cmd.Allocations) == 0 {
		return CreateWalletResult{}, httperr.NewIncorrectInputError(
			errors.New("wallet must have at least one allocation"),
			"missing-allocation",
		)
	}

	allocationDomainList, err := handler.buildAllocationListDomain(ctx, cmd.Allocations)
	if err != nil {
		return CreateWalletResult{}, httperr.NewIncorrectInputError(err, "fail-to-build-allocation")
	}

	walletDomain, err := handler.buildWalletDomain(cmd, allocationDomainList)
	if err != nil {
		return CreateWalletResult{}, httperr.NewIncorrectInputError(err, "fail-to-build-wallet")
	}

	if err = handler.walletRepo.Create(ctx, walletDomain); err != nil {
		return CreateWalletResult{}, httperr.NewUnknowError(err, "persist-wallet-failed")
	}

	return CreateWalletResult{
		WalletID: walletDomain.ID().String(),
	}, nil
}

func (handler *createWalletHandler) buildAllocationListDomain(ctx context.Context, allocationItems []CreateWalletAllocation) ([]*wallet.Allocation, error) {
	allocationDomainList := make([]*wallet.Allocation, 0, len(allocationItems))

	for _, item := range allocationItems {
		assetSourceID, err := assetsource.NewID(item.AssetSourceID)
		if err != nil {
			return nil, err
		}

		officeID, err := uuid.Parse(item.OfficeID)
		if err != nil {
			return nil, fmt.Errorf("invalid office id: %w", err)
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

		allocationDomain, err := wallet.NewAllocation(assetSourceID, amount, officeID)
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

	officeID, err := uuid.Parse(cmd.OfficeID)
	if err != nil {
		return nil, fmt.Errorf("invalid office id: %w", err)
	}

	walletDomain, err := wallet.NewWallet(
		cmd.Name,
		walletCurrency,
		cmd.IsStrictMode,
		officeID,
		allocationDomainList,
	)
	if err != nil {
		return nil, err
	}

	return walletDomain, nil
}
