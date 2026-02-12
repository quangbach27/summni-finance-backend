-- name: CreateFundProvider :exec
INSERT INTO finance.fund_providers (
    id,
    balance,
    currency,
    available_amount,
    version
) VALUES(
    $1, -- id
    $2, -- balance
    $3, -- currency
    $4, -- available_amount
    $5  -- version
);

-- name: GetFundProviderByID :one
SELECT 
    id,
    balance,
    currency,
    available_amount,
    version
FROM finance.fund_providers
WHERE id = $1;

-- name: GetFundProvidersByWalletID :many
SELECT
    fp.id,
    fp.balance,
    fp.currency,
    fp.available_amount,
    fp.version,
    fpa.allocated_amount
FROM finance.fund_providers fp
    LEFT JOIN finance.fund_provider_allocation fpa
        ON fp.id = fpa.fund_provider_id
WHERE fpa.wallet_id = $1;

-- name: UpdateFundProviderPartial :execrows
UPDATE finance.fund_providers
SET
    balance = COALESCE(sqlc.narg('balance'), balance),
    currency = COALESCE(sqlc.narg('currency'), currency),
    available_amount = COALESCE(sqlc.narg('available_amount'), available_amount),
    version = version + 1
WHERE id = sqlc.arg('id')
  AND version = sqlc.arg('version');