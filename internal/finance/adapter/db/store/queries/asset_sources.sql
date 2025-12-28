-- name: CreateAssetSource :exec
INSERT INTO finance.asset_sources (
    id, 
    owner_id, 
    balance, 
    source_type, 
    currency_code, 
    bank_name, 
    account_number,
    office_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

