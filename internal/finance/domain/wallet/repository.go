package wallet

import "context"

type Repository interface {
	GetByID(ctx context.Context, id ID) (*Wallet, error)
	Create(ctx context.Context, wallets *Wallet) error
}
