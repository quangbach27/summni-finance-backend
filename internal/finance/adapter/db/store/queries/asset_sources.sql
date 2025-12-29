-- name: CreateAssetSource :exec
INSERT INTO finance.asset_sources (
    id, 
    owner_id, 
    balance, 
    assetsource_name,
    source_type, 
    currency_code, 
    bank_name, 
    account_number,
    office_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetAssetSourceByID :one
SELECT 
    id,
    owner_id,
    assetsource_name,
    balance,
    source_type,
    currency_code,
    bank_name,
    account_number,
    office_id
FROM finance.asset_sources
WHERE id = $1;