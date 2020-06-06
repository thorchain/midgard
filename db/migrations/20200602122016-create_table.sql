
-- +migrate Up

CREATE TABLE blocks (
    time            TIMESTAMPTZ     NOT NULL,
    height          BIGINT          PRIMARY KEY,
)

CREATE TABLE events (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       PRIMARY KEY,
    height          BIGINT          NOT NULL,
    type            VARCHAR         NOT NULL,
    metadata        JSONB,
);

CREATE TYPE tx_direction as ENUM('in', 'out');
CREATE TABLE txs (
    time            TIMESTAMPTZ     NOT NULL,
    tx_hash         VARCHAR         PRIMARY KEY,
    event_id        BIGINT          NOT NULL,
    direction       tx_direction    NOT NULL,
    from_address    VARCHAR,
    to_address      VARCHAR,
    memo            VARCHAR,
);
CREATE INDEX direction_txs_idx ON txs (direction);
CREATE INDEX from_address_txs_idx ON txs (from_address);
CREATE INDEX to_address_txs_idx ON txs (to_address);

CREATE TABLE coins (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       PRIMARY KEY,
    tx_hash         VARCHAR         NOT NULL,
    chain           VARCHAR         NOT NULL,
    symbol          VARCHAR         NOT NULL,
    ticker          VARCHAR         NOT NULL,
    amount          BIGINT          NOT NULL,
);
CREATE INDEX tx_hash_coins_idx ON coins (tx_hash);
CREATE INDEX chain_coins_idx ON coins (chain);
CREATE INDEX symbol_coins_idx ON coins (symbol);
CREATE INDEX ticker_coins_idx ON coins (ticker);

CREATE TABLE pool_changes (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       PRIMARY KEY,
    event_id        BIGINT          NOT NULL,
    pool            VARCHAR         NOT NULL,
    asset_amount    BIGINT          NOT NULL,
    rune_amount     BIGINT          NOT NULL,
    units           BIGINT          NOT NULL,
    status          SMALLINT,
    tx_hash         VARCHAR,
);
CREATE INDEX event_id_pool_changes_idx ON pool_changes (event_id);
CREATE INDEX pool_pool_changes_idx ON pool_changes (pool);

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
SELECT create_hypertable('blocks', 'time');
SELECT create_hypertable('events', 'time');
SELECT create_hypertable('txs', 'time');
SELECT create_hypertable('coins', 'time');
SELECT create_hypertable('pool_changes', 'time');

-- +migrate Down

DROP TABLE blocks;
DROP TABLE events;
DROP TABLE txs;
DROP TABLE coins;
DROP TABLE pool_changes;

DROP TYPE tx_direction;