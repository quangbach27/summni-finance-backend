BEGIN;

-- 1. Create the finance schema
CREATE SCHEMA IF NOT EXISTS finance;

CREATE TABLE finance.fund_providers(
    id uuid PRIMARY KEY NOT NULL,
    
    balance bigint NOT NULL DEFAULT 0 
        CHECK (balance >= 0),

    currency char(3) NOT NULL,

    available_amount bigint NOT NULL DEFAULT 0 
        CHECK (available_amount >= 0),

    version int NOT NULL DEFAULT 0 
        CHECK (version >= 0),

    CHECK (available_amount <= balance)
);

CREATE TABLE finance.wallets (
    id uuid PRIMARY KEY NOT NULL,

    balance bigint NOT NULL 
        CHECK (balance >= 0),

    currency varchar(3) NOT NULL,

    version int NOT NULL
);

CREATE TABLE finance.fund_provider_allocation(
    fund_provider_id uuid NOT NULL,
    wallet_id uuid NOT NULL,

    allocated_amount bigint NOT NULL DEFAULT 0
        CHECK (allocated_amount >= 0),

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

COMMIT;