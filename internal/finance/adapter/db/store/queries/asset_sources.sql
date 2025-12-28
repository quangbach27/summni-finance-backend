-- name: ListAssetSources :many
SELECT id, owner_id, balance, source_type, currency_code, bank_name, account_number 
FROM finance.asset_sources
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateAssetsource :exec
INSERT INTO finance.asset_sources (
    id, 
    owner_id, 
    balance, 
    source_type, 
    currency_code, 
    bank_name, 
    account_number
)
VALUES ($1, $2, $3, $4, $5, $6, $7);

