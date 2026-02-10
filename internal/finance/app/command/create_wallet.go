package command

import "sumni-finance-backend/internal/common/cqrs"

type CreateWalletCmd struct {
}

type CreateWalletHandler cqrs.CommandHandler[CreateWalletCmd]

type createWalletHandler struct {
	
}
