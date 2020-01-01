
-- +migrate Up

CREATE TABLE genesis
(
    genesis_time TIMESTAMPTZ not null,
    primary key (genesis_time)
);

CREATE TABLE events (
    time timestamptz not null,
    id bigserial,
    event_id bigint not null,
    height bigint not null,
    type varchar not null,
    status varchar not null,
    to_address varchar,
    from_address varchar,
    pool varchar not null,
    rune_amount bigint,
    asset_amount bigint,
    stake_units bigint,
    swap_price_target bigint,
    swap_trade_slip real,
    swap_liquidity_fee bigint,
    primary key (id, time)
);

CREATE TYPE tx_direction as enum('in', 'out');
CREATE TABLE txs (
    time timestamptz not null,
    id bigserial,
    tx_hash varchar not null,
    event_id bigint not null,
    direction tx_direction not null,
    chain varchar ,
    from_address varchar,
    to_address varchar,
    memo varchar,
    gas_amount bigint,
    primary key(id, time)
);

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
SELECT create_hypertable('events', 'time');
SELECT create_hypertable('txs', 'time');
CREATE INDEX idx_events ON events (event_id, from_address, pool);
CREATE INDEX idx_txs ON txs (to_address, from_address, event_id);
-- +migrate Down

DROP TABLE genesis;
DROP TABLE events;

