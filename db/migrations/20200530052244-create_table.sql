
-- +migrate Up

CREATE TABLE events (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       NOT NULL,
    height          BIGINT          NOT NULL,
    type            VARCHAR         NOT NULL,
    status          VARCHAR,
    PRIMARY KEY (id, time)
);

CREATE TYPE event_type AS enum(
    'stake',
    'add',
    'unstake',
    'swap',
    'double_swap',
    'refund',
    'pool',
    'rewards',
    'gas',
    'fee',
    'slash',
    'errata',
    'outbound');
CREATE TYPE swap_type AS enum('buy', 'sell');
CREATE TYPE tx_direction AS enum('in', 'out');
CREATE TABLE pools_history (
    time            TIMESTAMPTZ     NOT NULL,
    id              BIGSERIAL       NOT NULL,
    type            event_type,
    event_id        BIGINT          NOT NULL,
    event_type      event_type      NOT NULL,
    pool            VARCHAR         NOT NULL,
    asset_amount    BIGINT          NOT NULL,
    asset_depth     BIGINT          NOT NULL,
    rune_amount     BIGINT          NOT NULL,
    rune_depth      BIGINT          NOT NULL,
    units           BIGINT,
    swap_type       swap_type,
    trade_slip      REAL,
    liquidity_fee   BIGINT,
    price_target    BIGINT,
    from_address    VARCHAR,
    to_address      VARCHAR,
    tx_hash         VARCHAR,
    tx_memo         VARCHAR,
    tx_direction    tx_direction,
    status          SMALLINT        NOT NULL,
    PRIMARY KEY (id, time)
);
CREATE INDEX pools_history_event_id_idx ON pools_history (event_id);
CREATE INDEX pools_history_event_type_idx ON pools_history USING hash (event_type);
CREATE INDEX pools_history_pool_idx ON pools_history (pool);
CREATE INDEX pools_history_swap_type ON pools_history USING hash (swap_type) WHERE swap_type IS NOT NULL;
CREATE INDEX pools_history_tx_hash ON pools_history USING hash (tx_hash) WHERE tx_hash IS NOT NULL;
CREATE INDEX pools_history_tx_direction ON pools_history USING hash (tx_direction) WHERE tx_direction IS NOT NULL;

CREATE TABLE pools (
    time              TIMESTAMPTZ              NOT NULL,
    pool              VARCHAR                  NOT NULL,
    height            BIGINT                   NOT NULL,
    asset_depth       BIGINT                   NOT NULL,
    asset_staked      BIGINT                   NOT NULL,
    asset_withdrawn   BIGINT                   NOT NULL,
    rune_depth        BIGINT                   NOT NULL,
    rune_staked       BIGINT                   NOT NULL,
    rune_withdrawn    BIGINT                   NOT NULL,
    units             BIGINT                   NOT NULL,
    status            SMALLINT                 NOT NULL,
    buy_volume        BIGINT                   NOT NULL,
    buy_slip_total    DOUBLE PRECISION         NOT NULL,
    buy_fee_total     BIGINT                   NOT NULL,
    buy_count         BIGINT                   NOT NULL,
    sell_volume       BIGINT                   NOT NULL,
    sell_slip_total   DOUBLE PRECISION         NOT NULL,
    sell_fee_total    BIGINT                   NOT NULL,
    sell_count        BIGINT                   NOT NULL,
    stakers_count     BIGINT                   NOT NULL,
    swappers_count    BIGINT                   NOT NULL,
    stake_count       BIGINT                   NOT NULL,
    withdraw_count    BIGINT                   NOT NULL,
    PRIMARY KEY (pool, time)
);

CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
SELECT create_hypertable('events', 'time');
SELECT create_hypertable('pools_history', 'time');
SELECT create_hypertable('pools', 'time');

-- +migrate Down

DROP TABLE events;
DROP TABLE pools_history;
DROP TABLE pools;

DROP TYPE event_type;
DROP TYPE swap_type;
DROP TYPE tx_direction;
