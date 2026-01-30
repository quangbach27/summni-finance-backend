-- Up Migration for sumni_finance_db

-- 1. Create the finance schema
CREATE SCHEMA IF NOT EXISTS finance;

CREATE TABLE finance.fund_providers(
    id uuid PRIMARY KEY NOT NULL,
    
    balance bigint NOT NULL DEFAULT 0,
    currency varchar(255) NOT NULL,
    fund_provider_type varchar(50) NOT NULL,

    -- Cash Details
    name varchar(255),

    -- Bank Details
    bank_name varchar(255),
    account_owner varchar(255),
    account_number varchar(255),

    version int DEFAULT 0
);

CREATE TABLE finance.wallets (
    id uuid PRIMARY KEY NOT NULL
)

CREATE TABLE finance.fund_provider_allocation(
    fund_provider_id uuid NOT NULL,
    wallet_id uuid NOT NULL,

    amount bigint NOT NULL DEFAULT 0,

    PRIMARY KEY (fund_provider_id, wallet_id),

    CONSTRAINT fk_fundprovider_allocation
        FOREIGN KEY (fund_provider_id)
            REFERENCES finance.fund_providers (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_wallet_allocation
        FOREIGN KEY (wallet_id)
            REFERENCES finance.wallets (id)
            ON DELETE CASCADE
);