package service

import (
	"sumni-finance-backend/internal/finance/app"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"
)

func NewApplication() app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateAssetSourceHandler: command.NewCreateAssetSourceHandler(),
		},
		Queries: app.Queries{
			GetAssetSourceHandler: query.NewGetAssetSoureHandler(),
		},
	}
}
