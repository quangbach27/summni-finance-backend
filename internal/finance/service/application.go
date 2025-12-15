package service

import (
	"context"
	common_db "sumni-finance-backend/internal/common/db"
	"sumni-finance-backend/internal/finance/adapter/db"
	"sumni-finance-backend/internal/finance/adapter/db/store"
	"sumni-finance-backend/internal/finance/app"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"
)

func NewApplication() app.Application {
	ctx := context.Background()
	connPool := common_db.MustNewPgConnectionPool(ctx)
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
