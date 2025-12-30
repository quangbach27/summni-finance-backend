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
    assetsource_id,
    wallet_id,
    amount
) VALUES ($1, $2, $3);

-- name: GetWalletsWithAllocationsByOfficeID :many
SELECT 
    w.id,
    w.wallet_name,
    w.currency_code,
    w.balance,
    w.is_strict_mode,
    w.office_id,
    a.id as assetsource_id,
    a.owner_id,
    a.assetsource_name,
    a.balance,
    a.source_type,
    a.currency_code,
    wa.amount
FROM finance.wallets w
    LEFT JOIN finance.wallets_allocation wa ON w.id = wa.wallet_id 
    LEFT JOIN finance.asset_sources a ON a.id = wa.assetsource_id
WHERE w.office_id = $1;