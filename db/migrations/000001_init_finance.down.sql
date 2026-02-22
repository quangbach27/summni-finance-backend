BEGIN;

-- 1. Drop dependent table first
DROP TABLE IF EXISTS finance.fund_provider_allocation;

-- 2. Drop base tables
DROP TABLE IF EXISTS finance.wallets;
DROP TABLE IF EXISTS finance.fund_providers;

-- 3. Drop schema (only if empty)
DROP SCHEMA IF EXISTS finance;

COMMIT;