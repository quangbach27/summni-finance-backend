-- Up Migration for sumni_finance_db

-- 1. Create the finance schema
CREATE SCHEMA IF NOT EXISTS finance;

-- Set the search path temporarily for this script to simplify table creation
SET search_path TO finance, public;

-- 2. CREATE TABLE assetsources
CREATE TABLE finance.asset_sources (
    id uuid PRIMARY KEY NOT NULL,
    owner_id uuid NOT NULL,
    
    balance bigint NOT NULL DEFAULT 0,
    source_type varchar(100) NOT NULL DEFAULT 'CASH',
    currency_code CHAR(3) NOT NULL DEFAULT 'VND',
    
    bank_name varchar(255),
    account_number varchar(255)
);

-- 3. CREATE TABLE wallets
CREATE TABLE finance.wallets(
    id uuid PRIMARY KEY NOT NULL,
    owner_id uuid NOT NULL,
    name varchar(255) NOT NULL,
    currency_code CHAR(3) NOT NULL,
    is_strict_mode boolean DEFAULT false
);

-- 4. CREATE TABLE assetsource_wallet (Junction Table)
CREATE TABLE finance.assetsource_wallet(
    asset_source_id uuid NOT NULL,
    wallet_id uuid NOT NULL,

    PRIMARY KEY (asset_source_id, wallet_id),

    CONSTRAINT fk_assetsource
        FOREIGN KEY (asset_source_id)
            REFERENCES finance.asset_sources (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_wallet
        FOREIGN KEY (wallet_id)
            REFERENCES finance.wallets (id)
            ON DELETE CASCADE
);