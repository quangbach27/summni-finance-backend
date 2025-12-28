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
	AssetSourceList []CreateAssetSourceItem
}

type CreateAssetSourceItem struct {
	Name          string
	OwnerID       uuid.UUID
	InitBalance   int64
	SourceType    string
	CurrencyCode  string
	BankName      string
	AccountNumber string
}

type CreateAssetSourceHandler cqrs.CommandHandler[CreateAssetSourceCmd]

type createAssetSourceHandler struct {
	assetsourceRepo assetsource.Repository
}

func NewCreateAssetSourceHandler(assetsourceRepo assetsource.Repository) CreateAssetSourceHandler {
	return cqrs.ApplyCommandDecorators(&createAssetSourceHandler{assetsourceRepo: assetsourceRepo})
}

func (h *createAssetSourceHandler) Handle(ctx context.Context, cmd CreateAssetSourceCmd) (err error) {
	assetSourceList := make([]*assetsource.AssetSource, 0, len(cmd.AssetSourceList))

	for _, cmd := range cmd.AssetSourceList {
		sourceType, err := assetsource.NewSourceTypeFromStr(cmd.SourceType)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "invalid-source-type")
		}

		currency, err := valueobject.NewCurrency(cmd.CurrencyCode)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "invalid-currency-code")
		}

		var assetSource *assetsource.AssetSource
		if sourceType.IsCash() {
			assetSource, err = assetsource.NewCashAssetSource(
				cmd.OwnerID,
				cmd.InitBalance,
				currency,
			)
			if err != nil {
				return httperr.NewError(err, "failed-to-create-cash-asset-source")
			}
		} else {
			assetSource, err = assetsource.NewBankAssetSource(
				cmd.OwnerID,
				cmd.InitBalance,
				currency,
				cmd.BankName,
				cmd.AccountNumber,
			)
			if err != nil {
				return httperr.NewError(err, "failed-to-create-bank-asset-source")
			}
		}

		assetSourceList = append(assetSourceList, assetSource)
	}

	if err = h.assetsourceRepo.Create(ctx, assetSourceList); err != nil {
		return httperr.NewUnknowError(err, "failed-to-store-asset-source-to-db")
	}

	return nil
}
