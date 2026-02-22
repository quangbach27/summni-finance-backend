package db

import (
	"context"
	"errors"
	"fmt"
	common_db "sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type walletRepo struct {
	queries            *store.Queries
	transactionManager *common_db.PgxTransactionManager
}

func NewWalletRepo(
	queries *store.Queries,
	transactionManager *common_db.PgxTransactionManager,
) (*walletRepo, error) {
	if queries == nil || transactionManager == nil {
		return nil, errors.New("missing dependencies")
	}

	return &walletRepo{
		queries:            queries,
		transactionManager: transactionManager,
	}, nil
}

func (r *walletRepo) GetByID(
	ctx context.Context,
	wID uuid.UUID,
) (*wallet.Wallet, error) {
	return nil, errors.New("walletRepo.GetByID not implemented")
}

func (r *walletRepo) GetByIDWithProviders(
	ctx context.Context,
	wID uuid.UUID,
	spec wallet.ProviderAllocationSpec,
) (*wallet.Wallet, error) {
	walletModel, err := r.queries.GetWalletByID(ctx, wID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve wallet '%s': %w", wID.String(), err)
	}

	providerModels, err := r.queries.GetFundProviderByWalletID(ctx, wID)
	if err != nil {
		return nil, err
	}

	filteredProviderAllocationsDomain := make([]wallet.ProviderAllocation, 0, len(providerModels))
	for _, model := range providerModels {
		fundProvider, err := fundprovider.UnmarshalFundProviderFromDatabase(
			model.ID,
			model.Balance,
			model.UnallocatedAmount,
			model.Currency,
			model.Version,
		)
		if err != nil {
			return nil, err
		}

		providerAllocation, err := wallet.NewProviderAllocation(fundProvider, model.WalletAllocatedAmount)
		if err != nil {
			return nil, err
		}

		if spec.IsSatisfiedBy(providerAllocation) {
			filteredProviderAllocationsDomain = append(filteredProviderAllocationsDomain, providerAllocation)
		}
	}

	walletDomain, err := wallet.UnmarshalWalletFromDatabase(
		walletModel.ID,
		walletModel.Balance,
		walletModel.Currency,
		walletModel.Version,
		filteredProviderAllocationsDomain...,
	)
	if err != nil {
		return nil, err
	}

	return walletDomain, nil
}

func (r *walletRepo) Create(ctx context.Context, wallet *wallet.Wallet) error {
	return r.queries.CreateWallet(ctx, store.CreateWalletParams{
		ID:       wallet.ID(),
		Balance:  wallet.Balance().Amount(),
		Currency: wallet.Currency().Code(),
		Version:  wallet.Version(),
	})
}

func (r *walletRepo) Update(
	ctx context.Context,
	wID uuid.UUID,
	spec wallet.ProviderAllocationSpec,
	updateFunc func(w *wallet.Wallet) error,
) (err error) {
	w, err := r.GetByIDWithProviders(ctx, wID, spec)
	if err != nil {
		return fmt.Errorf("failed to retrieve wallet :%w", err)
	}

	err = updateFunc(w)
	if err != nil {
		return err
	}

	return r.transactionManager.WithTx(ctx, func(tx pgx.Tx) error {
		qtx := r.queries.WithTx(tx)

		// update allocation
		for _, pa := range w.ProviderManager().ProviderAllocations() {
			err = qtx.UpsertFundProviderAllocation(
				ctx,
				store.UpsertFundProviderAllocationParams{
					FundProviderID:  pa.Provider().ID(),
					WalletID:        w.ID(),
					AllocatedAmount: pa.Allocated().Amount(),
				},
			)
			if err != nil {
				return fmt.Errorf("failed to update fund provider allocation: %w", err)
			}

			// update fundprovider
			rows, err := qtx.UpdateFundProviderPartial(ctx, store.UpdateFundProviderPartialParams{
				UnallocatedAmount: common_db.ToPgInt8(pa.Provider().UnallocatedBalance().Amount()),
				ID:                pa.Provider().ID(),
				Version:           pa.Provider().Version(),
			})
			if err != nil {
				return err
			}

			if rows == 0 {
				return fmt.Errorf("failed to update fund provider: %w", common_db.ErrConcurrentModification)
			}
		}

		// update wallet
		rows, err := qtx.UpdateWalletPartial(ctx, store.UpdateWalletPartialParams{
			ID:       w.ID(),
			Balance:  common_db.ToPgInt8(w.Balance().Amount()),
			Currency: common_db.ToPgText(w.Currency().Code()),
			Version:  w.Version(),
		})
		if err != nil {
			return err
		}
		if rows == 0 {
			return fmt.Errorf("failed to update wallet: %w", common_db.ErrConcurrentModification)
		}
		return nil
	})
}
