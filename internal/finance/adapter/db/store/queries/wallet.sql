-- name: CreateWallet :exec
INSERT INTO finance.wallets (
    id,
    name,
    currency_code,
    balance,
    is_strict_mode
)
VALUES ($1, $2, $3, $4, $5);

-- name: CreateWalletAssetSourceAssociateBatch :copyfrom
INSERT INTO finance.assetsource_wallet (
    asset_source_id,
    wallet_id
) VALUES ($1, $2);