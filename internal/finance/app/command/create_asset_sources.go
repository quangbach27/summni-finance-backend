package command

import (
	"context"
	"errors"
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
	OwnerID       string
	InitBalance   int64
	SourceType    string
	CurrencyCode  string
	BankName      string
	AccountNumber string
	OfficeID      string
}

type CreateAssetSourceHandler cqrs.CommandHandler[CreateAssetSourceCmd]

type createAssetSourceHandler struct {
	assetsourceRepo assetsource.Repository
}

func NewCreateAssetSourceHandler(assetsourceRepo assetsource.Repository) CreateAssetSourceHandler {
	return cqrs.ApplyCommandDecorators(&createAssetSourceHandler{assetsourceRepo: assetsourceRepo})
}

func (h *createAssetSourceHandler) Handle(ctx context.Context, cmd CreateAssetSourceCmd) (err error) {
	if len(cmd.AssetSourceList) == 0 {
		return httperr.NewIncorrectInputError(errors.New("asset source list cannot be empty"), "asset-source-list-is-empty")
	}

	assetSourceList := make([]*assetsource.AssetSource, 0, len(cmd.AssetSourceList))

	for _, asCmd := range cmd.AssetSourceList {
		sourceType, err := assetsource.NewSourceTypeFromStr(asCmd.SourceType)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "invalid-source-type")
		}

		currency, err := valueobject.NewCurrency(asCmd.CurrencyCode)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "invalid-currency-code")
		}

		ownerID, err := uuid.Parse(asCmd.OwnerID)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "fail-to-parse-owner-id")
		}

		officeID, err := uuid.Parse(asCmd.OfficeID)
		if err != nil {
			return httperr.NewIncorrectInputError(err, "fail-to-parse-office-id")
		}

		var assetSource *assetsource.AssetSource
		if sourceType.IsCash() {
			assetSource, err = assetsource.NewCashAssetSource(
				ownerID,
				asCmd.InitBalance,
				currency,
				officeID,
			)
			if err != nil {
				return httperr.NewIncorrectInputError(err, "failed-to-create-cash-asset-source")
			}
		} else {
			assetSource, err = assetsource.NewBankAssetSource(
				ownerID,
				asCmd.InitBalance,
				currency,
				asCmd.BankName,
				asCmd.AccountNumber,
				officeID,
			)
			if err != nil {
				return httperr.NewIncorrectInputError(err, "failed-to-create-bank-asset-source")
			}
		}

		assetSourceList = append(assetSourceList, assetSource)
	}

	if err = h.assetsourceRepo.Create(ctx, assetSourceList); err != nil {
		return httperr.NewUnknowError(err, "failed-to-save-asset-source-to-db")
	}

	return nil
}
