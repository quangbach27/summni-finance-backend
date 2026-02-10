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