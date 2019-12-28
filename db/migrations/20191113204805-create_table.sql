
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
    tx_in_memo varchar,
    tx_out_memo varchar,
    tx_in_hash varchar,
    tx_out_hash varchar,
    tx_in_gas_chain varchar,
    tx_out_gas_chain varchar,
    tx_in_gas_amount bigint,
    tx_out_gas_amount bigint,
    primary key (id, time)
);

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
SELECT create_hypertable('events', 'time');
CREATE INDEX idx_events ON events (event_id, from_address, pool);

-- +migrate Down

DROP TABLE genesis;
DROP TABLE events;

