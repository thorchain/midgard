
-- +migrate Up
-- CREATE TABLE events (
--     time        TIMESTAMPTZ       NOT NULL,
--     id integer not null,
--     type varchar not null,
--     height integer not null,
--     in_tx jsonb not null,
--     out_tx jsonb not null,
--     gas jsonb not null,
--     event jsonb not null,
--     status varchar not null,
--     primary key (time, id)
-- );


CREATE TABLE events (
    time        TIMESTAMPTZ       not null,
    id integer not null,
    height integer not null,
    type varchar not null,
    status varchar not null,
    primary key (time, id)
);

CREATE TABLE stakes (
    time        TIMESTAMPTZ       NOT NULL,
    event_id integer references events(id),
    chain integer not null,
    symbol varchar not null,
    ticker varchar not null,
    units integer,
    primary key (time, event_id)
);
--
-- CREATE TYPE tx_direction as enum('in', 'out');
-- CREATE TABLE txs (
--     time        TIMESTAMPTZ       NOT NULL,
--     hash varchar not null,
--     event_id integer references events(id),
--     direction tx_direction not null,
--     chain varchar,
--     from_address varchar,
--     to_address varchar,
--     memo varchar,
--     primary key (time, event_id)
-- );
--
-- CREATE TABLE coins (
--     time        TIMESTAMPTZ       NOT NULL,
--     tx_id integer references txs(id),
--     event_id integer references events(id),
--     chain varchar not null,
--     symbol varchar not null,
--     ticker varchar not null,
--     amount integer not null,
--     primary key (time, tx_id, event_id)
-- );
--
-- CREATE TABLE gas_coins (
--     time        TIMESTAMPTZ       NOT NULL,
--     event_id integer references events(id),
--     chain varchar not null,
--     symbol varchar not null,
--     ticker varchar not null,
--     amount integer not null,
--     primary key (time, event_id)
-- );

SELECT create_hypertable('events', 'time');
SELECT create_hypertable('stakes', 'time');
-- SELECT create_hypertable('txs', 'time', 'event_id');
-- SELECT create_hypertable('coins', 'time', 'tx_id', 'event_id');
-- SELECT create_hypertable('gas_coins', 'time', 'event_id');

-- +migrate Down

DROP TABLE events;

