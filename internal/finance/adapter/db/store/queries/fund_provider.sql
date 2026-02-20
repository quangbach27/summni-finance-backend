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
