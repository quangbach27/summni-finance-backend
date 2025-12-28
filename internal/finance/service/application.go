package service

import (
	"sumni-finance-backend/internal/finance/adapter/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/app"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewApplication(connPool *pgxpool.Pool) app.Application {
	queries := store.New(connPool)

	assetRepo := db.NewAssetsourceRepo(connPool, queries)

	return app.Application{
		Commands: app.Commands{
			CreateAssetSourceHandler: command.NewCreateAssetSourceHandler(assetRepo),
		},
		Queries: app.Queries{
			GetAssetSourceHandler: query.NewGetAssetSoureHandler(),
		},
	}
}
