-- name: CreateAccountingPeriod :exec
INSERT INTO finance.accounting_periods (
    id,
    year_month,
    start_date,
    interval,
    end_time,
    wallet_opening_balance,
    total_debit,
    total_credit,
    wallet_closing_balance,
    version,
    wallet_id,
    status
) VALUES (
    $1, -- id
    $2, -- year_month
    $3, -- start_date
    $4, -- interval
    $5, -- end_time
    $6, -- wallet_opening_balance
    $7, -- total_debit
    $8, -- total_credit
    $9, -- wallet_closing_balance
    $10, -- version
    $11, -- wallet_id
    $12 -- status
);

-- name: UpdateAccountingPerid :execrows
UPDATE finance.accounting_periods ap
SET
    total_debit = sqlc.arg(total_debit),
    total_credit = sqlc.arg(total_credit),
    wallet_closing_balance = sqlc.arg(closing_balance),
    status = sqlc.arg(status),
    version = version + 1
WHERE ap.id = sqlc.arg(id)
    AND ap.version = sqlc.arg(version);

-- name: GetAccountingPeriodsByYearMonthAndWalletID :one
SELECT
    id,
    year_month,
    start_date,
    interval,
    end_time,
    wallet_opening_balance,
    total_debit,
    total_credit,
    wallet_closing_balance,
    status,
    version
FROM
    finance.accounting_periods
WHERE wallet_id = $1 
    AND year_month = $2;

-- name: GetWalletWithAccountingPeriod :one
SELECT
    w.id             AS wallet_id,
    w.name           AS wallet_name,
    w.balance        AS wallet_balance,
    w.currency       AS wallet_currency,
    w.version        AS wallet_version,
    ap.id            AS period_id,
    ap.year_month    AS period_year_month,
    ap.start_date    AS period_start_date,
    ap.interval      AS period_interval,
    ap.end_time      AS period_end_time,
    ap.wallet_opening_balance,
    ap.total_debit,
    ap.total_credit,
    ap.wallet_closing_balance,
    ap.status        AS period_status,
    ap.version       AS period_version
FROM finance.wallets w
LEFT JOIN finance.accounting_periods ap
    ON ap.wallet_id = w.id
    AND ap.year_month = $2
WHERE w.id = $1;

-- name: BulkInsertTransactionRecords :copyfrom
INSERT INTO finance.transaction_records (
    id,
    transaction_no,
    transaction_type,
    amount,
    wallet_balance,
    wallet_id,
    fp_id,
    fp_balance,
    accounting_periods_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
);
