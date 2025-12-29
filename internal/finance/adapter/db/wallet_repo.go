package db

import (
	"context"
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type walletRepository struct {
	pool    *pgxpool.Pool
	queries *store.Queries
}

func NewWalletRepository(connPool *pgxpool.Pool, queries *store.Queries) (*walletRepository, error) {
	if connPool == nil {
		return nil, errors.New("missing connection pool")
	}

	if queries == nil {
		return nil, errors.New("missing queried")
	}

	return &walletRepository{
		pool:    connPool,
		queries: queries,
	}, nil
}

func (repo *walletRepository) GetByID(ctx context.Context, id wallet.ID) (*wallet.Wallet, error) {
	return nil, nil
}

func (repo *walletRepository) Create(ctx context.Context, wallet *wallet.Wallet) (err error) {
	balance, err := wallet.TotalBalance()
	if err != nil {
		return fmt.Errorf("fail to calculate balance: %w", err)
	}

	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		err = db.FinishTransaction(ctx, tx, err)
	}()

	txQueries := repo.queries.WithTx(tx)

	arg := store.CreateWalletParams{
		ID:           uuid.UUID(wallet.ID()),
		Name:         wallet.Name(),
		CurrencyCode: wallet.Currency().Code(),
		Balance:      balance.Amount(),
		IsStrictMode: wallet.IsStrictMode(),
	}

	err = txQueries.CreateWallet(ctx, arg)
	if err != nil {
		return err
	}

	associateParams := make([]store.CreateWalletAssetSourceAssociateBatchParams, 0, len(wallet.Allocations()))
	for _, allocation := range wallet.Allocations() {
		associateParams = append(associateParams, store.CreateWalletAssetSourceAssociateBatchParams{
			AssetSourceID: uuid.UUID(allocation.AssetSourceID()),
			WalletID:      uuid.UUID(wallet.ID()),
		})
	}

	_, err = txQueries.CreateWalletAssetSourceAssociateBatch(ctx, associateParams)
	if err != nil {
		return err
	}

	return nil
}
