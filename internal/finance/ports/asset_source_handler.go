package ports

import (
	"encoding/json"
	"errors"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type CreateAssetSourceRequest struct {
	AssetSources []CreateAssetSourceItem `json:"assetSource"`
}

type CreateAssetSourceItem struct {
	Name          string `json:"name"`
	OwnerID       string `json:"ownerId"`
	InitBalance   int64  `json:"initBalance"`
	SourceType    string `json:"sourceType"`
	CurrencyCode  string `json:"currencyCode"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
}

func (h *FinanceHandler) CreateAssetSource(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateAssetSourceRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		httperr.BadRequest("failed-to-parse-json", err, w, r)
		return
	}

	if len(req.AssetSources) == 0 {
		httperr.BadRequest("assetSource-empty", errors.New("assetSource is required"), w, r)
		return
	}

	cmd := command.CreateAssetSourceCmd{
		AssetSourceList: make([]command.CreateAssetSourceItem, len(req.AssetSources)),
	}
	for i, as := range req.AssetSources {
		ownerID, err := uuid.Parse(as.OwnerID)
		if err != nil {
			httperr.BadRequest("invalid-ownerID", err, w, r)
			return
		}

		cmd.AssetSourceList[i] = command.CreateAssetSourceItem{
			Name:          as.Name,
			OwnerID:       ownerID,
			InitBalance:   as.InitBalance,
			SourceType:    as.SourceType,
			CurrencyCode:  as.CurrencyCode,
			BankName:      as.BankName,
			AccountNumber: as.AccountNumber,
		}
	}

	err := h.app.Commands.CreateAssetSourceHandler.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]any{
		"status": "created",
	})
}

func (h *FinanceHandler) GetAssetSources(w http.ResponseWriter, r *http.Request) {
	cmd := query.GetAssetSourceCmd{}

	_, err := h.app.Queries.GetAssetSourceHandler.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}
}
