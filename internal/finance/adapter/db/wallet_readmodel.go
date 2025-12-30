package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/app/query"

	"github.com/google/uuid"
)

type walletReadModel struct {
	queries *store.Queries
}

func NewWalletReadModel(queries *store.Queries) (*walletReadModel, error) {
	if queries == nil {
		return nil, errors.New("missing store queries")
	}

	return &walletReadModel{queries: queries}, nil
}

func (r *walletReadModel) GetAllWalletsWithAllocations(ctx context.Context, officeID uuid.UUID) ([]query.Wallet, error) {
	walletModelsWithAllocation, err := r.queries.GetWalletsWithAllocationsByOfficeID(ctx, officeID)
	if err != nil {
		return nil, err
	}
	slog.Info(fmt.Sprintf("models: %d", len(walletModelsWithAllocation)))

	walletAppMap := make(map[uuid.UUID]query.Wallet, len(walletModelsWithAllocation))
	for _, walletModel := range walletModelsWithAllocation {
		slog.Info(walletModel.ID.String())

		wallet, exist := walletAppMap[walletModel.ID]
		if !exist {
			slog.Info("Does not exist. Create New UUID in map")
			walletApp := r.toWalletApp(walletModel)
			walletAppMap[walletModel.ID] = walletApp
			continue
		}

		slog.Info("wallet already exist. Append allocation")
		wallet.Allocations = append(wallet.Allocations, r.toAllocationApp(walletModel))
	}
	slog.Info(fmt.Sprintf("map : %d", len(walletAppMap)))

	walletApps := make([]query.Wallet, 0, len(walletAppMap))
	for _, v := range walletAppMap {
		walletApps = append(walletApps, v)
	}

	slog.Info(fmt.Sprintf("walletAPP: %d", len(walletApps)))

	return walletApps, nil
}

func (r *walletReadModel) toWalletApp(model store.GetWalletsWithAllocationsByOfficeIDRow) query.Wallet {
	return query.Wallet{
		Name:         model.WalletName,
		Balance:      model.Balance,
		CurrencyCode: model.CurrencyCode,
		IsStrictMode: model.IsStrictMode,
		Allocations:  []query.Allocation{r.toAllocationApp(model)},
	}
}

func (r *walletReadModel) toAllocationApp(model store.GetWalletsWithAllocationsByOfficeIDRow) query.Allocation {
	return query.Allocation{
		AssetSourceID:   model.AssetsourceID.String(),
		AssetSourceName: model.AssetsourceName.String,
		Amount:          model.Amount.Int64,
	}
}
