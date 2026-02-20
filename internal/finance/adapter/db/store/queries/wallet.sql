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
