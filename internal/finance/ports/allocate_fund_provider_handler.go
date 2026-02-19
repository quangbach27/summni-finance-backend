package ports

import (
	"encoding/json"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/common/server/response"
	"sumni-finance-backend/internal/finance/app/command"
)

// Allocate fund providers to a wallet
// (PUT /v1/wallet/allocate-fund-provider)
func (hs HttpServer) AllocateFundProvider(w http.ResponseWriter, r *http.Request) {
	var req AllocateFundProviderRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		httperr.BadRequest("failed-to-parse-json", err, w, r)
		return
	}

	fundProviders := make([]command.FundProviderCmd, 0, len(req.Providers))
	for _, provider := range req.Providers {
		fundProviders = append(fundProviders, command.FundProviderCmd{
			ID:        provider.ProviderID,
			Allocated: provider.Allocated,
		})
	}

	err := hs.application.Commands.AllocateFundProvider.Handle(r.Context(), command.AllocateFundProviderCmd{
		WalletID:      req.WalletId,
		FundProviders: fundProviders,
	})

	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	response.WriteJSON(w, r, http.StatusOK, nil, nil)
}
