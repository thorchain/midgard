
-- +migrate Up

CREATE TABLE events (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       NOT NULL,
    height          BIGINT          NOT NULL,
    type            VARCHAR         NOT NULL,
    status          VARCHAR,
    PRIMARY KEY (id, time)
);

CREATE TABLE pools_history (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       NOT NULL,
    height          BIGINT          NOT NULL,
    event_id        BIGINT          NOT NULL,
    event_type      VARCHAR         NOT NULL,
    pool            VARCHAR         NOT NULL,
    asset_amount    BIGINT          NOT NULL,
    asset_depth     BIGINT          NOT NULL,
    rune_amount     BIGINT          NOT NULL,
    rune_depth      BIGINT          NOT NULL,
    units           BIGINT,
    status          SMALLINT        NOT NULL,
    meta            JSONB,
    PRIMARY KEY (id, time)
);
CREATE INDEX pools_history_event_id_idx ON pools_history (event_id);
CREATE INDEX pools_history_event_type_idx ON pools_history (event_type);
CREATE INDEX pools_history_pool_idx ON pools_history (pool);

CREATE TABLE swaps (
    time        TIMESTAMPTZ       NOT NULL,
    id SERIAL,
    event_id bigint not null,
    from_address varchar not null,
    to_address varchar not null,
    pool varchar not null,
    price_target bigint,
    trade_slip real,
    liquidity_fee bigint,
    runeAmt bigint,
    assetAmt bigint,
    primary key (id, time)
);
CREATE INDEX idx_swaps ON swaps (from_address, pool);

CREATE TYPE tx_direction as enum('in', 'out');
CREATE TABLE txs (
    time        TIMESTAMPTZ       NOT NULL,
    id SERIAL,
    tx_hash varchar not null,
    event_id bigint not null,
    direction tx_direction not null,
    chain varchar,
    from_address varchar,
    to_address varchar,
    memo varchar,
    primary key (id, time, event_id)
);

CREATE TABLE coins (
    time        TIMESTAMPTZ       NOT NULL,
    id SERIAL,
    tx_hash varchar not null,
    event_id bigint not null,
    chain varchar not null,
    symbol varchar not null,
    ticker varchar not null,
    amount bigint not null,
    primary key (id, time, event_id)
);

CREATE TABLE stakers (
     time                TIMESTAMPTZ NOT NULL,
     id                  SERIAL,
     rune_address        VARCHAR NOT NULL,
     asset_address       VARCHAR NOT NULL,
     pool                VARCHAR NOT NULL,
     unit                BIGINT NOT NULL,
     asset_staked        BIGINT NOT NULL,
     rune_staked         BIGINT NOT NULL,
     asset_withdrawn     BIGINT NOT NULL,
     rune_withdrawn      BIGINT NOT NULL,
     height_last_staked  BIGINT,
     height_first_staked BIGINT,
     PRIMARY KEY (rune_address, asset_address, pool, time)
);

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
SELECT create_hypertable('events', 'time');
SELECT create_hypertable('pools_history', 'time');
SELECT create_hypertable('swaps', 'time');
SELECT create_hypertable('txs', 'time');
SELECT create_hypertable('coins', 'time');
SELECT create_hypertable('stakers', 'time');

-- +migrate Down

DROP TABLE events;
DROP TABLE pools_history;
DROP TABLE swaps;
DROP TABLE txs;
DROP TABLE coins;
DROP TABLE stakers;

DROP TYPE tx_direction;
