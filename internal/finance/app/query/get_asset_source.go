package query

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/common/cqrs"
)

type GetAssetSourceCmd struct{}

type GetAssetSourceHandler cqrs.QueryHandler[GetAssetSourceCmd, AssetSource]

type getAssetSourceHandler struct{}

func NewGetAssetSoureHandler() GetAssetSourceHandler {
	return cqrs.ApplyQueryDecorator(&getAssetSourceHandler{})
}

func (h *getAssetSourceHandler) Handle(ctx context.Context, cmd GetAssetSourceCmd) (AssetSource, error) {

	return AssetSource{}, errors.New("here is my error")
}
