-- name: CreateFundProvider :exec
INSERT INTO finance.fund_providers (
    id,
    balance,
    currency,
    version
) VALUES(
    $1, -- id
    $2, -- balance
    $3, -- currency
    $4  -- version
);

-- name: GetFundProviderByID :one
SELECT 
    id,
    balance,
    currency,
    version,
    fp.balance - COALESCE(SUM(fpa.allocated_amount), 0) as available_amount
FROM finance.fund_providers fp
    LEFT JOIN finance.fund_provider_allocation fpa
        ON fp.id = fpa.fund_provider_id
WHERE id = $1
GROUP BY fp.id;

