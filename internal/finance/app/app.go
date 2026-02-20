package app

import (
	"sumni-finance-backend/internal/finance/adapter/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/app/command"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateFundProvider command.CreateFundProviderHandler
	CreateWallet       command.CreateWalletHandler
}

type Queries struct {
}

func NewApplication(pgPool *pgxpool.Pool) (Application, error) {
	queries := store.New(pgPool)

	walletRepo, err := db.NewWalletRepo(queries, pgPool)
	if err != nil {
		return Application{}, err
	}

	fundProviderRepo, err := db.NewFundProviderRepo(queries, pgPool)
	if err != nil {
		return Application{}, err
	}

	return Application{
		Commands: Commands{
			CreateFundProvider: command.NewCreateFundProviderHandler(fundProviderRepo),
			CreateWallet:       command.NewCreateWalletHandler(walletRepo),
		},
		Queries: Queries{},
	}, nil
}
