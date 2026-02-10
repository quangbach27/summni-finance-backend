package app

import (
	"sumni-finance-backend/internal/common/cqrs"
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

func NewApplication(
	pgxPool *pgxpool.Pool,
) Application {
	quries := store.New(pgxPool)

	fundProviderRepo, err := db.NewFundProviderRepository(quries)
	if err != nil {
		panic(err)
	}

	walletProviderRepo, err := db.NewWalletRepository(quries, pgxPool)
	if err != nil {
		panic(err)
	}

	return Application{
		Commands: Commands{
			CreateFundProvider: cqrs.ApplyCommandDecorators(command.NewCreateFundProviderHandler(fundProviderRepo)),
			CreateWallet:       cqrs.ApplyCommandDecorators(command.NewCreateWalletHandler(walletProviderRepo, fundProviderRepo)),
		},
		Queries: Queries{},
	}
}
