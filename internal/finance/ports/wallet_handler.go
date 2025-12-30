package ports

import (
	"encoding/json"
	"errors"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"
	"sumni-finance-backend/internal/finance/app/query"

	"github.com/go-chi/render"
)

type CreateWalletRequest struct {
	Name         string              `json:"name"`
	CurrencyCode string              `json:"currencyCode"`
	IsStrictMode bool                `json:"isStrictMode"`
	OfficeID     string              `json:"officeId"`
	Allocations  []AllocationRequest `json:"allocations"`
}

type AllocationRequest struct {
	AssetSourceID string `json:"assetSourceId"`
	Amount        int64  `json:"amount"`
	OfficeID      string `json:"officeId"`
}

func (h *financeHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req CreateWalletRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		httperr.BadRequest("failed-to-parse-json", err, w, r)
		return
	}

	if len(req.Allocations) == 0 {
		httperr.BadRequest("missing-allocations", errors.New("missing allocations"), w, r)
		return
	}

	allocationItemList := make([]command.CreateWalletAllocation, 0, len(req.Allocations))
	for _, alloc := range req.Allocations {
		allocationItem := command.CreateWalletAllocation{
			AssetSourceID: alloc.AssetSourceID,
			Amount:        alloc.Amount,
			OfficeID:      alloc.OfficeID,
		}

		allocationItemList = append(allocationItemList, allocationItem)
	}

	cmd := command.CreateWalletCmd{
		Name:         req.Name,
		CurrencyCode: req.CurrencyCode,
		IsStrictMode: req.IsStrictMode,
		OfficeID:     req.OfficeID,
		Allocations:  allocationItemList,
	}

	result, err := h.app.Commands.CreateWalletHandler.Handle(r.Context(), cmd)
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

func (h *financeHandler) GetAllWallets(w http.ResponseWriter, r *http.Request) {
	rawOfficeID := r.URL.Query().Get("office_id")

	cmd := query.GetWalletListCmd{
		OfficeID: rawOfficeID,
	}

	walletApps, err := h.app.Queries.GetWalletListHandler.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, envelop{"wallets": walletApps})
}
