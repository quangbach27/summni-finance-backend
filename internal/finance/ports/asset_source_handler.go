package ports

import (
	"encoding/json"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"

	"github.com/go-chi/render"
)

type CreateAssetSourceRequest struct {
	Name          string `json:"name"`
	OwnerID       string `json:"ownerId"`
	InitBalance   int64  `json:"initBalance"`
	SourceType    string `json:"sourceType"`
	CurrencyCode  string `json:"currencyCode"`
	BankName      string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	OfficeID      string `json:"officeId"`
}

func (h *financeHandler) CreateAssetSources(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateAssetSourceRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		httperr.BadRequest("failed-to-parse-json", err, w, r)
		return
	}

	cmd := command.CreateAssetSourceCmd{
		Name:          req.Name,
		OwnerID:       req.OwnerID,
		InitBalance:   req.InitBalance,
		SourceType:    req.SourceType,
		CurrencyCode:  req.CurrencyCode,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		OfficeID:      req.OfficeID,
	}

	result, err := h.app.Commands.CreateAssetSourceHandler.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, envelop{
		"status": "created",
		"data":   result,
	})
}

func (h *financeHandler) GetAssetSources(w http.ResponseWriter, r *http.Request) {
	cmd := query.GetAssetSourceCmd{}

	_, err := h.app.Queries.GetAssetSourceHandler.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}
}
