package db

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/db"
	commondb "sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrConcurrentModification = errors.New("concurrent modification detected")

type walletRepository struct {
	queries *store.Queries
	pgPool  *pgxpool.Pool
}

func NewWalletRepository(queries *store.Queries, pgPool *pgxpool.Pool) (*walletRepository, error) {
	if queries == nil {
		return nil, errors.New("queries is required")
	}

	if pgPool == nil {
		return nil, errors.New("pgPool is required")
	}

	return &walletRepository{
		queries: queries,
		pgPool:  pgPool,
	}, nil
}

func (r *walletRepository) GetByID(ctx context.Context, wID uuid.UUID) (*wallet.Wallet, error) {
	// Get wallet and its allocation
	walletModel, err := r.queries.GetWalletByID(ctx, wID)
	if err != nil {
		return nil, err
	}

	// Get fundProvider with list of walletID
	fundProviderModels, err := r.queries.GetFundProvidersByWalletID(ctx, wID)
	if err != nil {
		return nil, err
	}

	providerAllocations := make([]wallet.ProviderAllocation, 0, len(fundProviderModels))

	for _, model := range fundProviderModels {
		currency, err := valueobject.NewCurrency(model.Currency)
		if err != nil {
			return nil, err
		}

		balance, err := valueobject.NewMoney(model.Balance, currency)
		if err != nil {
			return nil, err
		}
		availableAmount, err := valueobject.NewMoney(model.AvailableAmount, currency)

		fundProvider, err := fundprovider.UnmarshallFundProviderFromDatabase(
			model.ID,
			balance,
			availableAmount,
			model.Version,
		)
		if err != nil {
			return nil, err
		}

		allocated, err := valueobject.NewMoney(model.AvailableAmount, currency)

		providerAllocation, err := wallet.NewProviderAllocation(fundProvider, allocated)
		if err != nil {
			return nil, err
		}

		providerAllocations = append(providerAllocations, providerAllocation)
	}

	currency := valueobject.Currency{}
	balance, err := valueobject.NewMoney(walletModel.Balance, currency)
	if err != nil {
		return nil, err
	}

	walletDomain, err := wallet.UnmarshalWalletFromDatabase(
		walletModel.ID,
		balance,
		walletModel.Version,
		providerAllocations...,
	)

	return walletDomain, nil
}

func (r *walletRepository) Create(ctx context.Context, wallet *wallet.Wallet) error {
	tx, err := r.pgPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = commondb.FinishTransaction(ctx, tx, err)
	}()

	// Use queries with transaction
	qtx := r.queries.WithTx(tx)

	// 1. Create Wallet
	err = qtx.CreateWallet(ctx, store.CreateWalletParams{
		ID:       wallet.ID(),
		Balance:  wallet.Balance().Amount(),
		Currency: wallet.Balance().Currency().Code(),
		Version:  wallet.Version(),
	})
	if err != nil {
		return err
	}

	providerAllocations := wallet.ProviderManager().GetFundProviderAllocations()
	if len(providerAllocations) == 0 {
		return nil
	}

	createFundProviderAllocationParams := make([]store.CreatFundProviderAllocationParams, 0, len(providerAllocations))
	updateFundProviderParams := []store.UpdateFundProviderPartialParams{}

	for _, allocation := range providerAllocations {
		createFundProviderAllocationParams = append(createFundProviderAllocationParams, store.CreatFundProviderAllocationParams{
			WalletID:        wallet.ID(),
			FundProviderID:  allocation.FundProvider().ID(),
			AllocatedAmount: allocation.Allocated().Amount(),
		})

		updateFundProviderParams = append(updateFundProviderParams, store.UpdateFundProviderPartialParams{
			AvailableAmount: db.ToPgInt8(allocation.FundProvider().AvailableAmountForAllocation().Amount()),
			Version:         allocation.FundProvider().Verions(),
		})
	}

	// 2. Batch create provider allocations
	rows, err := qtx.CreatFundProviderAllocation(ctx, createFundProviderAllocationParams)
	if err != nil {
		return fmt.Errorf("failed to create fund provider allocation: %v", err)
	}
	if rows != int64(len(createFundProviderAllocationParams)) {
		return errors.New("failed to partial fund_provider_allocation insert")
	}

	// 3. Update FundProvider
	for _, params := range updateFundProviderParams {
		rows, err = qtx.UpdateFundProviderPartial(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to update fund provider: %v", err)
		}

		if rows == 0 {
			return fmt.Errorf("failed to update fund provider: %v", ErrConcurrentModification)
		}
	}

	return nil
}

func (r *walletRepository) Update(ctx context.Context, wallet *wallet.Wallet) error {
	tx, err := r.pgPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = commondb.FinishTransaction(ctx, tx, err)
	}()

	qtx := r.queries.WithTx(tx)

	rows, err := qtx.UpdateWalletPartial(ctx, store.UpdateWalletPartialParams{
		ID:       wallet.ID(),
		Balance:  db.ToPgInt8(wallet.Balance().Amount()),
		Currency: db.ToPgText(wallet.Balance().Currency().Code()),
		Version:  wallet.Version(),
	})
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrConcurrentModification
	}

	// for _, allocation := range wallet.ProviderManager().GetFundProviderAllocations() {

	// }
	return nil
}
