package wallet

import "context"

type Repository interface {
	GetWalletAllocations(ctx context.Context, id ID) (*Wallet, error)
}
