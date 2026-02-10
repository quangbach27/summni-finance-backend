package ports

import (
	"encoding/json"
	"net/http"
	"sumni-finance-backend/internal/common/server/httperr"
	"sumni-finance-backend/internal/finance/app/command"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// Create a new fund provider
// (POST /v1/fund-provider)
func (hs HttpServer) CreateFundProvider(w http.ResponseWriter, r *http.Request) {
	var fundProviderReq CreateFundProviderRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&fundProviderReq); err != nil {
		httperr.BadRequest("failed-to-parse-json", err, w, r)
		return
	}

	cmd := command.CreateFundProviderCmd{
		Balance:  fundProviderReq.Balance,
		Currency: fundProviderReq.Currency,
	}

	err := hs.application.Commands.CreateFundProvider.Handle(r.Context(), cmd)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	requestID := middleware.GetReqID(r.Context())

	response := map[string]string{
		"request_id": requestID,
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}
