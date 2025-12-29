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
	assetSource, err := h.buildAssetSource(cmd)
	if err != nil {
		return httperr.NewIncorrectInputError(err, "fail-to-create-asset-source")
	}

	if err = h.assetsourceRepo.Create(ctx, assetSource); err != nil {
		return httperr.NewUnknowError(err, "failed-to-save-asset-source")
	}

	return nil
}

func (h *createAssetSourceHandler) buildAssetSource(cmd CreateAssetSourceCmd) (*assetsource.AssetSource, error) {
	sourceType, err := assetsource.NewSourceTypeFromStr(cmd.SourceType)
	if err != nil {
		return nil, err
	}

	currency, err := valueobject.NewCurrency(cmd.CurrencyCode)
	if err != nil {
		return nil, err
	}

	ownerID, err := uuid.Parse(cmd.OwnerID)
	if err != nil {
		return nil, err
	}

	officeID, err := uuid.Parse(cmd.OfficeID)
	if err != nil {
		return nil, err
	}

	if sourceType.IsCash() {
		cashAssetSource, err := assetsource.NewCashAssetSource(
			ownerID,
			cmd.Name,
			cmd.InitBalance,
			currency,
			officeID,
		)
		if err != nil {
			return nil, err
		}

		return cashAssetSource, nil
	}

	bankAssetSource, err := assetsource.NewBankAssetSource(
		ownerID,
		cmd.Name,
		cmd.InitBalance,
		currency,
		cmd.BankName,
		cmd.AccountNumber,
		officeID,
	)
	if err != nil {
		return nil, err
	}

	return bankAssetSource, nil
}
