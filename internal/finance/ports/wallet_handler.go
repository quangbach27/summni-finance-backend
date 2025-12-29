package ports

import (
	"encoding/json"
	"errors"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"

	"github.com/go-chi/render"
)

type CreateWalletRequest struct {
	Name         string              `json:"name"`
	CurrencyCode string              `json:"currencyCode"`
	IsStrictMode bool                `json:"isStrictMode"`
	Allocations  []AllocationRequest `json:"allocations"`
}

type AllocationRequest struct {
	AssetSourceID string `json:"assetSourceId"`
	Amount        int64  `json:"amount"`
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
		httperr.BadRequest("missing-allocations", errors.New("Missing allocations"), w, r)
		return
	}

	allocationItemList := make([]command.AllocationItem, 0, len(req.Allocations))
	for _, allocation := range req.Allocations {
		allocationItem := command.AllocationItem{
			AssetSourceID: allocation.AssetSourceID,
			Amount:        allocation.Amount,
		}

		allocationItemList = append(allocationItemList, allocationItem)
	}

	cmd := command.CreateWalletCmd{
		Name:         req.Name,
		CurrencyCode: req.CurrencyCode,
		IsStrictMode: req.IsStrictMode,
		Allocations:  allocationItemList,
	}

	if err := h.app.Commands.CreateWalletHandler.Handle(r.Context(), cmd); err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]any{
		"status": "created",
	})
}
