-- name: CreateWallet :exec
INSERT INTO finance.wallets (
    id,
    name,
    balance,
    currency,
    version
) VALUES (
    $1, -- id
    $2, -- name
    $3, -- balance
    $4, -- currency
    $5 -- version
);

-- name: GetWalletByID :one
SELECT 
    id,
    name,
    balance,
    currency,
    version
FROM finance.wallets
WHERE id = $1;

-- name: UpdateWalletPartial :execrows
UPDATE finance.wallets
SET
    name = COALESCE(sqlc.narg(name), name),
    balance  = COALESCE(sqlc.narg(balance), balance),
    currency = COALESCE(sqlc.narg(currency), currency),
    version  = version + 1
WHERE id = sqlc.arg(id)
  AND version = sqlc.arg(version);

-- name: UpsertFundProviderAllocation :exec
INSERT INTO finance.fund_provider_allocation (
    fund_provider_id,
    wallet_id,
    allocated_amount
) VALUES (
    $1,
    $2,
    $3
)
ON CONFLICT (fund_provider_id, wallet_id)
DO UPDATE
SET allocated_amount = EXCLUDED.allocated_amount;
