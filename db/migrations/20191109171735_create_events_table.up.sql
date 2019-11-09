CREATE TABLE IF NOT EXISTS events (
                            time        TIMESTAMPTZ       NOT NULL
);

SELECT create_hypertable('events', 'time');
