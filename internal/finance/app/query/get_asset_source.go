package query

import (
	"context"
	"sumni-finance-backend/internal/common/cqrs"
)

type GetAssetSourceCmd struct{}

type GetAssetSourceHandler cqrs.QueryHandler[GetAssetSourceCmd, AssetSource]

type getAssetSourceHandler struct{}

func NewGetAssetSoureHandler() GetAssetSourceHandler {
	return cqrs.ApplyQueryDecorator(&getAssetSourceHandler{})
}

func (h *getAssetSourceHandler) Handle(ctx context.Context, cmd GetAssetSourceCmd) (AssetSource, error) {
	return AssetSource{
		Name:       "techcombank",
		SourceType: "bank",
	}, nil
}
