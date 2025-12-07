package handler

import (
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"

	"github.com/go-chi/render"
)

func (h *FinanceHandler) CreateAssetSource(w http.ResponseWriter, r *http.Request) {
	cmd := command.CreateAssetSourceCmd{}

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
