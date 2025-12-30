package query

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"

	"github.com/google/uuid"
)

type GetWalletListCmd struct {
	OfficeID string
}

type Wallet struct {
	Name         string       `json:"name"`
	Balance      int64        `json:"balance"`
	CurrencyCode string       `json:"currencyCode"`
	IsStrictMode bool         `json:"isStrictMode"`
	Allocations  []Allocation `json:"allocations"`
}

type Allocation struct {
	AssetSourceID   string `json:"assetSourceId"`
	AssetSourceName string `json:"assetSourceName"`
	Amount          int64  `json:"amount"`
}

type GetWalletListHandler cqrs.QueryHandler[GetWalletListCmd, []Wallet]

type WalletReadModel interface {
	GetAllWalletsWithAllocations(ctx context.Context, officeID uuid.UUID) ([]Wallet, error)
}

type getWalletListHandler struct {
	walletReadModel WalletReadModel
}

func NewGetWalletListHandler(walletReadModel WalletReadModel) GetWalletListHandler {
	return cqrs.ApplyQueryDecorator(&getWalletListHandler{
		walletReadModel: walletReadModel,
	})
}

func (h *getWalletListHandler) Handle(ctx context.Context, cmd GetWalletListCmd) ([]Wallet, error) {
	if cmd.OfficeID == "" {
		return nil, httperr.NewIncorrectInputError(errors.New("officeID is required"), "missing-office-id")
	}

	officeID, err := uuid.Parse(cmd.OfficeID)
	if err != nil {
		return nil, httperr.NewIncorrectInputError(err, "invalid-office-id")
	}

	wallets, err := h.walletReadModel.GetAllWalletsWithAllocations(ctx, officeID)
	if err != nil {
		return nil, httperr.NewIncorrectInputError(err, "fail-to-retrieve-wallets")
	}

	return wallets, nil
}
