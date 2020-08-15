
-- +migrate Up

CREATE VIEW pool_changes_5_min WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s', timescaledb.materialized_only = true)
AS
SELECT pool, time_bucket('5 min', time) AS time,
    SUM(asset_amount) AS asset_changes,
    last(asset_depth, id) AS asset_depth,
    SUM(CASE WHEN event_type = 'stake' THEN asset_amount ELSE 0 END) AS asset_staked,
    SUM(CASE WHEN event_type = 'unstake' AND asset_amount < 0 THEN -asset_amount ELSE 0 END) AS asset_withdrawn,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    last(rune_depth, id) AS rune_depth,
    SUM(CASE WHEN event_type = 'stake' THEN rune_amount ELSE 0 END) AS rune_staked,
    SUM(CASE WHEN event_type = 'unstake' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS rune_withdrawn,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(units) AS units_changes,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY pool, time_bucket('5 min', time);

CREATE VIEW pool_changes_hourly WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s', timescaledb.materialized_only = true)
AS
SELECT pool, time_bucket('1 hour', time) AS time,
    SUM(asset_amount) AS asset_changes,
    last(asset_depth, id) AS asset_depth,
    SUM(CASE WHEN event_type = 'stake' THEN asset_amount ELSE 0 END) AS asset_staked,
    SUM(CASE WHEN event_type = 'unstake' AND asset_amount < 0 THEN -asset_amount ELSE 0 END) AS asset_withdrawn,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    last(rune_depth, id) AS rune_depth,
    SUM(CASE WHEN event_type = 'stake' THEN rune_amount ELSE 0 END) AS rune_staked,
    SUM(CASE WHEN event_type = 'unstake' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS rune_withdrawn,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(units) AS units_changes,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY pool, time_bucket('1 hour', time);

CREATE VIEW pool_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s', timescaledb.materialized_only = true)
AS
SELECT pool, time_bucket('1 day', time) AS time,
    SUM(asset_amount) AS asset_changes,
    last(asset_depth, id) AS asset_depth,
    SUM(CASE WHEN event_type = 'stake' THEN asset_amount ELSE 0 END) AS asset_staked,
    SUM(CASE WHEN event_type = 'unstake' AND asset_amount < 0 THEN -asset_amount ELSE 0 END) AS asset_withdrawn,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    last(rune_depth, id) AS rune_depth,
    SUM(CASE WHEN event_type = 'stake' THEN rune_amount ELSE 0 END) AS rune_staked,
    SUM(CASE WHEN event_type = 'unstake' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS rune_withdrawn,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(units) AS units_changes,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY pool, time_bucket('1 day', time);

CREATE VIEW total_volume_changes_5_min WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s')
AS
SELECT time_bucket('5 minute', time) AS time,
    SUM(CASE WHEN rune_amount > 0 AND event_type = 'swap' THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(CASE WHEN rune_amount < 0 AND event_type = 'swap' THEN -rune_amount ELSE 0 END) AS sell_volume
FROM pools_history
GROUP BY time_bucket('5 minute', time);

CREATE VIEW total_volume_changes_hourly WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s')
AS
SELECT time_bucket('1 hour', time) AS time,
    SUM(CASE WHEN rune_amount > 0 AND event_type = 'swap' THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(CASE WHEN rune_amount < 0 AND event_type = 'swap' THEN -rune_amount ELSE 0 END) AS sell_volume
FROM pools_history
GROUP BY time_bucket('1 hour', time);

CREATE VIEW total_volume_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_interval = '3s')
AS
SELECT time_bucket('1 day', time) AS time,
    SUM(CASE WHEN rune_amount > 0 AND event_type = 'swap' THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(CASE WHEN rune_amount < 0 AND event_type = 'swap' THEN -rune_amount ELSE 0 END) AS sell_volume
FROM pools_history
GROUP BY time_bucket('1 day', time);

-- +migrate Down

DROP VIEW pool_changes_5_min CASCADE;
DROP VIEW pool_changes_hourly CASCADE;
DROP VIEW pool_changes_daily CASCADE;
DROP VIEW total_volume_changes_5_min CASCADE;
DROP VIEW total_volume_changes_hourly CASCADE;
DROP VIEW total_volume_changes_daily CASCADE;