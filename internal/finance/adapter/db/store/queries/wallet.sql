-- name: CreateWallet :exec
INSERT INTO finance.wallets (
    id,
    wallet_name,
    currency_code,
    balance,
    is_strict_mode,
    office_id
)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateWalletsAllocationBatch :copyfrom
INSERT INTO finance.wallets_allocation (
    asset_source_id,
    wallet_id,
    amount
) VALUES ($1, $2, $3);
