
-- +migrate Up
CREATE TABLE events (
    time        TIMESTAMPTZ       NOT NULL,
    id integer not null ,
    state varchar not null,
    height integer not null,
    type varchar not null,
    in_hash varchar not null,
    out_hash varchar,
    in_memo varchar,
    out_memo varchar,
    from_address varchar,
    to_address varchar,
    from_coin varchar ,
    to_coin varchar,
    gas varchar,
    primary key (time, id)
);

SELECT create_hypertable('events', 'time', 'id');

-- +migrate Down

DROP TABLE events;
