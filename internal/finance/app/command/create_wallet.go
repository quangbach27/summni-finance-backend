package command

import (
	"context"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/domain/wallet"
)

type CreateWalletCmd struct {
	CurrencyCode string
}

type CreateWalletHandler cqrs.CommandHandler[CreateWalletCmd]

type createWalletHandler struct {
	walletRepo wallet.Repository
}

func NewCreateWalletHandler(walletRepo wallet.Repository) *createWalletHandler {
	return &createWalletHandler{
		walletRepo: walletRepo,
	}
}

func (h *createWalletHandler) Handle(ctx context.Context, cmd CreateWalletCmd) error {
	walletDomain, err := wallet.NewWallet(cmd.CurrencyCode)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-currency-code")
	}

	return h.walletRepo.Create(ctx, walletDomain)
}
