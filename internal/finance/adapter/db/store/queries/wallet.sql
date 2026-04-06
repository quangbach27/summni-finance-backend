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

-- name: UpdateWalletBalance :execrows
UPDATE finance.wallets
SET
    balance = sqlc.arg(balance),
    version = version + 1
WHERE id = sqlc.arg(id)
    AND version = sqlc.arg(version);

-- name: BulkInsertFundAllocations :copyfrom
INSERT INTO finance.fund_provider_allocations (
    fp_id,
    wallet_id,
    allocated_amount
) VALUES ($1, $2, $3);

-- name: BatchUpdateFundAllocations :exec
UPDATE finance.fund_provider_allocations
SET allocated_amount = data.allocated_amount
FROM (
    SELECT
        unnest(sqlc.arg(fp_ids)::uuid[])        AS fp_id,
        unnest(sqlc.arg(wallet_ids)::uuid[])    AS wallet_id,
        unnest(sqlc.arg(allocated_amounts)::bigint[]) AS allocated_amount
) AS data
WHERE finance.fund_provider_allocations.fp_id = data.fp_id
    AND finance.fund_provider_allocations.wallet_id = data.wallet_id;
