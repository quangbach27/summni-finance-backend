package handler

import (
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/query"
)

func (h *FinanceHandler) GetAssetSources(w http.ResponseWriter, r *http.Request) {
	cmd := query.GetAssetSourceCmd{}

	_, err := h.app.Queries.GetAssetSourceHandler.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
	}
}
