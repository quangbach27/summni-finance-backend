-- name: CreateWallet :exec
INSERT INTO finance.wallets (
    id,
    balance,
    currency,
    version
) VALUES (
    $1, -- id
    $2, -- balance
    $3, -- currency
    $4 -- version
);

-- name: CreatFundProviderAllocation :copyfrom
INSERT INTO finance.fund_provider_allocation (
    fund_provider_id,
    wallet_id,
    allocated_amount
) VALUES (
    $1, -- fund_provider_id
    $2, -- wallet_id
    $3 -- allocated_amount
);

-- name: GetWalletByID :one
SELECT
    id, 
    balance,
    currency,
    version
FROM finance.wallets
WHERE id = $1;

-- name: UpdateWalletPartial :execrows
UPDATE finance.wallets
SET
    balance = COALESCE(sqlc.narg('balance'), balance),
    currency = COALESCE(sqlc.narg('currency'), currency),
    version = version + 1
WHERE id = sqlc.arg('id')
  AND version = sqlc.arg('version');

-- name: UpdateFundProviderAllocation :exec
UPDATE finance.fund_provider_allocation
SET
    allocated_amount = $1
WHERE fund_provider_id = $2 AND wallet_id = $3;
