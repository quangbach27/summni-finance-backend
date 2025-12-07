package app

import (
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateAssetSourceHandler command.CreateAssetSourceHandler
}

type Queries struct {
	GetAssetSourceHandler query.GetAssetSourceHandler
}
