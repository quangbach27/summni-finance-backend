package ports

import (
	"encoding/json"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/server/response"
	"sumni-finance-backend/internal/finance/app/command"
)

// Create a new wallet
// (POST /v1/wallet)
func (hs HttpServer) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var createWalletReq CreateWalletRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&createWalletReq); err != nil {
		httperr.BadRequest("failed-to-parse-json", err, w, r)
		return
	}

	cmd := command.CreateWalletCmd{
		Currency:    createWalletReq.Currency,
		Allocations: make([]command.CreateWalletCmdAllocation, 0, len(createWalletReq.Allocations)),
	}

	for _, allocationReq := range createWalletReq.Allocations {
		cmd.Allocations = append(cmd.Allocations, command.CreateWalletCmdAllocation{
			ProviderID: allocationReq.ProviderID,
			Allocated:  allocationReq.Allocated,
		})
	}

	err := hs.application.Commands.CreateWallet.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	response.WriteJSON(w, r, http.StatusCreated, nil, nil)
}
