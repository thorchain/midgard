
-- +migrate Up

CREATE VIEW pool_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s')
AS
SELECT pool, event_type, time_bucket('1 day', time) AS time,
    SUM(CASE WHEN asset_amount > 0 THEN asset_amount ELSE 0 END) AS pos_asset_changes,
    SUM(CASE WHEN asset_amount < 0 THEN asset_amount ELSE 0 END) AS neg_asset_changes,
    SUM(CASE WHEN rune_amount > 0 THEN rune_amount ELSE 0 END) AS pos_rune_changes,
    SUM(CASE WHEN rune_amount < 0 THEN rune_amount ELSE 0 END) AS neg_rune_changes,
    SUM(units) AS units_changes
FROM pools_history
GROUP BY pool, event_type, time_bucket('1 day', time);

CREATE VIEW total_volume_changes_5_min WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s')
AS
SELECT time_bucket('5 minute', time) AS time,
    SUM(CASE WHEN rune_amount > 0 THEN rune_amount ELSE 0 END) AS pos_changes,
    SUM(CASE WHEN rune_amount < 0 THEN rune_amount ELSE 0 END) AS neg_changes,
FROM pools_history
GROUP BY time_bucket('5 minute', time);

CREATE VIEW total_volume_changes_hourly WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s')
AS
SELECT time_bucket('1 hour', time) AS time,
    SUM(CASE WHEN rune_amount > 0 THEN rune_amount ELSE 0 END) AS pos_changes,
    SUM(CASE WHEN rune_amount < 0 THEN rune_amount ELSE 0 END) AS neg_changes,
FROM pools_history
GROUP BY time_bucket('1 hour', time);

CREATE VIEW total_volume_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s')
AS
SELECT time_bucket('1 day', time) AS time,
    SUM(CASE WHEN rune_amount > 0 THEN rune_amount ELSE 0 END) AS pos_changes,
    SUM(CASE WHEN rune_amount < 0 THEN rune_amount ELSE 0 END) AS neg_changes,
FROM pools_history
GROUP BY time_bucket('1 day', time);

-- +migrate Down

DROP VIEW pool_changes_daily CASCADE;
DROP VIEW total_volume_changes_5_min CASCADE;
DROP VIEW total_volume_changes_hourly CASCADE;
DROP VIEW total_volume_changes_daily CASCADE;