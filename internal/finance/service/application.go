package service

import (
	"sumni-finance-backend/internal/finance/adapter/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/app"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewApplication(connPool *pgxpool.Pool) (app.Application, error) {
	queries := store.New(connPool)

	assetSourceRepo, err := db.NewAssetsourceRepo(connPool, queries)
	if err != nil {
		return app.Application{}, err
	}

	walletRepo, err := db.NewWalletRepository(connPool, queries)
	if err != nil {
		return app.Application{}, err
	}

	return app.Application{
		Commands: app.Commands{
			CreateAssetSourceHandler: command.NewCreateAssetSourceHandler(assetSourceRepo),
			CreateWalletHandler:      command.NewCreateWalletHandler(walletRepo, assetSourceRepo),
		},
		Queries: app.Queries{
			GetAssetSourceHandler: query.NewGetAssetSoureHandler(),
		},
	}, nil
}
