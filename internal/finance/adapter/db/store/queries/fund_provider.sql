-- name: CreateFundProvider :exec
INSERT INTO finance.fund_providers (
    id,
    balance,
    currency,
    unallocated_amount,
    version
) VALUES(
    $1, -- id
    $2, -- balance
    $3, -- currency
    $4, -- unallocated_amount
    $5  -- version
);


-- name: GetFundProviderByWalletID :many
SELECT 
    fp.id,
    fp.balance,
    fp.currency,
    fp.unallocated_amount,
    fp.version,
    fpa.allocated_amount AS wallet_allocated_amount
FROM finance.fund_providers fp
    INNER JOIN finance.fund_provider_allocation fpa
        ON fp.id = fpa.fund_provider_id
            AND fpa.wallet_id = $1;

-- name: UpdateFundProviderPartial :execrows
UPDATE finance.fund_providers
SET
    balance = COALESCE(sqlc.narg(balance), balance),
    unallocated_amount = COALESCE(sqlc.narg(unallocated_amount), unallocated_amount),
    currency = COALESCE(sqlc.narg(currency), currency),
    version = version + 1
WHERE id = sqlc.arg(id)
  AND version = sqlc.arg(version);