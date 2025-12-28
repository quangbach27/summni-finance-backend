package command

import (
	"context"
	"sumni-finance-backend/internal/common/cqrs"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/assetsource"

	"github.com/google/uuid"
)

type CreateAssetSourceCmd struct {
	Name          string
	OwnerID       string
	InitBalance   int64
	SourceType    string
	CurrencyCode  string
	BankName      string
	AccountNumber string
	OfficeID      string
}

type CreateAssetSourceItem struct {
}

type CreateAssetSourceHandler cqrs.CommandHandler[CreateAssetSourceCmd]

type createAssetSourceHandler struct {
	assetsourceRepo assetsource.Repository
}

func NewCreateAssetSourceHandler(assetsourceRepo assetsource.Repository) CreateAssetSourceHandler {
	return cqrs.ApplyCommandDecorators(&createAssetSourceHandler{assetsourceRepo: assetsourceRepo})
}

func (h *createAssetSourceHandler) Handle(ctx context.Context, cmd CreateAssetSourceCmd) (err error) {
	sourceType, err := assetsource.NewSourceTypeFromStr(cmd.SourceType)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-source-type")
	}

	currency, err := valueobject.NewCurrency(cmd.CurrencyCode)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "invalid-currency-code")
	}

	ownerID, err := uuid.Parse(cmd.OwnerID)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "fail-to-parse-owner-id")
	}

	officeID, err := uuid.Parse(cmd.OfficeID)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "fail-to-parse-office-id")
	}

	var assetSource *assetsource.AssetSource
	if sourceType.IsCash() {
		assetSource, err = assetsource.NewCashAssetSource(
			ownerID,
			cmd.InitBalance,
			currency,
			officeID,
		)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "failed-to-create-cash-asset-source")
		}
	} else {
		assetSource, err = assetsource.NewBankAssetSource(
			ownerID,
			cmd.InitBalance,
			currency,
			cmd.BankName,
			cmd.AccountNumber,
			officeID,
		)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "failed-to-create-bank-asset-source")
		}
	}

	if err = h.assetsourceRepo.Create(ctx, assetSource); err != nil {
		return httperr.NewUnknowError(err, "failed-to-save-asset-source-to-db")
	}

	return nil
}
