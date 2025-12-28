package query

import (
	"context"
	"sumni-finance-backend/internal/common/cqrs"
)

type GetAssetSourceCmd struct {
	OfficeID string
}
type AssetSource struct {
	Name          string `json:"name"`
	SourceType    string `json:"sourceType"`
	Balance       int64  `json:"balance"`
	Currency      string `json:"currency"`
	BankName      string `json:"bankName,omitempty"`
	AccountNumber string `json:"accountNumber,omitempty"`
}

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
