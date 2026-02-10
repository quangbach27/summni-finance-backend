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
    version
FROM finance.fund_providers
WHERE id = $1;

