package command

import (
	"context"
	"errors"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/domain/ledger"
	"sumni-finance-backend/internal/finance/domain/wallet"

	"github.com/google/uuid"
)

type OpenAccountingPeriodCmd struct {
	WalletID uuid.UUID
	Year     int
	Month    int
}

type OpenAccountingPeriodHandler cqrs.CommandHandler[OpenAccountingPeriodCmd]

type openAccountingPeriodHandler struct {
	walletRepo wallet.Repository
	ledgerRepo ledger.Repository
}

func NewOpenAccountingPeriodHandler(walletRepo wallet.Repository, ledgerRepo ledger.Repository) OpenAccountingPeriodHandler {
	return &openAccountingPeriodHandler{
		walletRepo: walletRepo,
		ledgerRepo: ledgerRepo,
	}
}

func (h *openAccountingPeriodHandler) Handle(ctx context.Context, cmd OpenAccountingPeriodCmd) error {
	newYearMonth, err := ledger.NewYearMonth(cmd.Month, cmd.Year)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "failed-to-create-year-month")
	}

	w, err := h.walletRepo.GetByIDWithAccountingPeriod(ctx, cmd.WalletID, newYearMonth)
	if err != nil {
		return httperr.NewUnknowError(err, "failed-to-retrieve-wallet")
	}

	if err := w.OpenAccountingPeriod(newYearMonth); err != nil {
		return httperr.NewIncorrectInputError(err, "failed-to-open-accounting-period")
	}

	ap, exist := w.LedgerManager().FindAccountingPeriod(newYearMonth)
	if !exist {
		return httperr.NewUnknowError(
			errors.New("accounting period is successful opened in domain but not found in wallet domain"),
			"failed-to-open-accounting-period",
		)
	}

	if err = h.ledgerRepo.CreateAccountingPeriod(ctx, w.ID(), ap); err != nil {
		return err
	}

	return nil
}
