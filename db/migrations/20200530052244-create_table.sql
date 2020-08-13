
-- +migrate Up

CREATE TABLE events (
    time                TIMESTAMPTZ     NOT NULL,
    id                  BIGSERIAL       NOT NULL,
    height              BIGINT          NOT NULL,
    type                VARCHAR         NOT NULL,
    status              VARCHAR,
    swap_price_target   BIGINT,
    PRIMARY KEY (id, time)
);

CREATE TYPE swap_type as enum('buy', 'sell');
CREATE TABLE pools_history (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       NOT NULL,
    event_id        BIGINT          NOT NULL,
    event_type      VARCHAR         NOT NULL,
    pool            VARCHAR         NOT NULL,
    asset_amount    BIGINT          NOT NULL,
    asset_depth     BIGINT          NOT NULL,
    rune_amount     BIGINT          NOT NULL,
    rune_depth      BIGINT          NOT NULL,
    units           BIGINT,
    swap_type       swap_type,
    trade_slip      REAL,
    liquidity_fee   BIGINT,
    status          SMALLINT        NOT NULL,
    PRIMARY KEY (id, time)
);
CREATE INDEX pools_history_event_id_idx ON pools_history (event_id);
CREATE INDEX pools_history_event_type_idx ON pools_history (event_type);
CREATE INDEX pools_history_pool_idx ON pools_history (pool);
CREATE INDEX pools_history_swap_type ON pools_history (swap_type);

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

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
SELECT create_hypertable('events', 'time');
SELECT create_hypertable('pools_history', 'time');
SELECT create_hypertable('txs', 'time');
SELECT create_hypertable('coins', 'time');

-- +migrate Down

DROP TABLE events;
DROP TABLE pools_history;
DROP TABLE txs;
DROP TABLE coins;

DROP TYPE swap_type;
DROP TYPE tx_direction;
