
-- +migrate Up

CREATE TYPE event_type AS enum(
    'stake',
    'add',
    'unstake',
    'swap',
    'refund',
    'pool',
    'rewards',
    'gas',
    'fee',
    'slash',
    'errata',
    'outbound',
    'bond');

CREATE TYPE event_status AS enum('unknown', 'success');

CREATE TABLE events (
    time            TIMESTAMP       NOT NULL,
    height          BIGINT          NOT NULL,
    id              BIGSERIAL       NOT NULL,
    type            event_type      NOT NULL,
    event_id        BIGINT          NOT NULL,
    event_type      event_type      NOT NULL,
    event_status    event_status    NOT NULL,
    pool            VARCHAR,
    asset_amount    BIGINT,
    rune_amount     BIGINT,
    meta            JSONB,
    from_address    VARCHAR,
    to_address      VARCHAR,
    tx_hash         VARCHAR,
    tx_memo         VARCHAR,
    PRIMARY KEY (id, time)
);

CREATE TABLE pools_history (
    time              TIMESTAMP                NOT NULL,
    height            BIGINT                   NOT NULL,
    pool              VARCHAR                  NOT NULL,
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

CREATE TABLE pools (
    asset   VARCHAR PRIMARY KEY
);

CREATE TABLE stats_history (
    time                TIMESTAMP   NOT NULL,
    height              BIGINT      NOT NULL,
    total_users         BIGINT      NOT NULL,
    total_txs           BIGINT      NOT NULL,
    total_volume        BIGINT      NOT NULL,
    total_staked        BIGINT      NOT NULL,
    total_earned        BIGINT      NOT NULL,
    rune_depth          BIGINT      NOT NULL,
    pools_count         BIGINT      NOT NULL,
    buys_count          BIGINT      NOT NULL,
    sells_count         BIGINT      NOT NULL,
    stakes_count        BIGINT      NOT NULL,
    withdrawns_count    BIGINT      NOT NULL,
    PRIMARY KEY (time)
);

CREATE TABLE stakers (
    address             VARCHAR     NOT NULL,
    pool                VARCHAR     NOT NULL,
    units               BIGINT      NOT NULL,
    asset_staked        BIGINT      NOT NULL,
    asset_withdrawn     BIGINT      NOT NULL,
    rune_staked         BIGINT      NOT NULL,
    rune_withdrawn      BIGINT      NOT NULL,
    first_stake_at      TIMESTAMP   NOT NULL,
    last_stake_at       TIMESTAMP   NOT NULL,
    last_withdrawn_at   TIMESTAMP,
    PRIMARY KEY (address, pool)
);

SELECT create_hypertable('events', 'time');
SELECT create_hypertable('pools_history', 'time');
SELECT create_hypertable('stats_history', 'time');

-- +migrate Down

DROP TABLE events;
DROP TABLE pools_history;
DROP TABLE pools;
DROP TABLE stats_history;
DROP TABLE stakers;

DROP TYPE event_type;
DROP TYPE event_status;
