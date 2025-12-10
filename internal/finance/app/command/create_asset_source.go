package command

import (
	"context"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/logs"
)

type CreateAssetSourceCmd struct{}

type CreateAssetSourceHandler cqrs.CommandHandler[CreateAssetSourceCmd]

type createAssetSourceHandler struct {
}

func NewCreateAssetSourceHandler() CreateAssetSourceHandler {
	return cqrs.ApplyCommandDecorators(&createAssetSourceHandler{})
}

func (h *createAssetSourceHandler) Handle(ctx context.Context, cmd CreateAssetSourceCmd) error {
	logs.FromContext(ctx).Info("create asset source.")
	return nil
}
