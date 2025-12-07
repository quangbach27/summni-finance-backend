package handler

import (
	"sumni-finance-backend/internal/finance/app"
	"sumni-finance-backend/internal/finance/ports"
)

type FinanceHandler struct {
	app app.Application
}

func NewFinanceServerInterface(app app.Application) ports.FinanceServerInterface {
	return &FinanceHandler{
		app: app,
	}
}
