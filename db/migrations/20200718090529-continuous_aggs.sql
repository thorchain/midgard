
-- +migrate Up

CREATE VIEW pool_changes_5_min WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT pool, time_bucket('5 min', time) AS time,
    SUM(asset_amount) AS asset_changes,
    last(asset_depth, id) AS asset_depth,
    SUM(CASE WHEN event_type = 'stake' THEN asset_amount ELSE 0 END) AS asset_staked,
    SUM(CASE WHEN event_type = 'unstake' AND asset_amount < 0 THEN -asset_amount ELSE 0 END) AS asset_withdrawn,
    SUM(CASE WHEN event_type = 'add' THEN asset_amount ELSE 0 END) AS asset_added,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    last(rune_depth, id) AS rune_depth,
    SUM(CASE WHEN event_type = 'stake' THEN rune_amount ELSE 0 END) AS rune_staked,
    SUM(CASE WHEN event_type = 'unstake' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS rune_withdrawn,
	SUM(CASE WHEN event_type = 'add' THEN rune_amount ELSE 0 END) AS rune_added,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(units) AS units_changes,
    SUM(CASE WHEN event_type = 'rewards' THEN rune_amount ELSE 0 END) AS reward,
	SUM(CASE WHEN event_type = 'gas' THEN -asset_amount ELSE 0 END) AS gas_used,
	SUM(CASE WHEN event_type = 'gas' THEN rune_amount ELSE 0 END) AS gas_replenished,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY pool, time_bucket('5 min', time);

CREATE VIEW pool_changes_hourly WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT pool, time_bucket('1 hour', time) AS time,
    SUM(asset_amount) AS asset_changes,
    last(asset_depth, id) AS asset_depth,
    SUM(CASE WHEN event_type = 'stake' THEN asset_amount ELSE 0 END) AS asset_staked,
    SUM(CASE WHEN event_type = 'unstake' AND asset_amount < 0 THEN -asset_amount ELSE 0 END) AS asset_withdrawn,
    SUM(CASE WHEN event_type = 'add' THEN asset_amount ELSE 0 END) AS asset_added,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    last(rune_depth, id) AS rune_depth,
    SUM(CASE WHEN event_type = 'stake' THEN rune_amount ELSE 0 END) AS rune_staked,
    SUM(CASE WHEN event_type = 'unstake' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS rune_withdrawn,
	SUM(CASE WHEN event_type = 'add' THEN rune_amount ELSE 0 END) AS rune_added,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(units) AS units_changes,
    SUM(CASE WHEN event_type = 'rewards' THEN rune_amount ELSE 0 END) AS reward,
	SUM(CASE WHEN event_type = 'gas' THEN -asset_amount ELSE 0 END) AS gas_used,
	SUM(CASE WHEN event_type = 'gas' THEN rune_amount ELSE 0 END) AS gas_replenished,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY pool, time_bucket('1 hour', time);

CREATE VIEW pool_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT pool, time_bucket('1 day', time) AS time,
    SUM(asset_amount) AS asset_changes,
    last(asset_depth, id) AS asset_depth,
    SUM(CASE WHEN event_type = 'stake' THEN asset_amount ELSE 0 END) AS asset_staked,
    SUM(CASE WHEN event_type = 'unstake' AND asset_amount < 0 THEN -asset_amount ELSE 0 END) AS asset_withdrawn,
    SUM(CASE WHEN event_type = 'add' THEN asset_amount ELSE 0 END) AS asset_added,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    last(rune_depth, id) AS rune_depth,
    SUM(CASE WHEN event_type = 'stake' THEN rune_amount ELSE 0 END) AS rune_staked,
    SUM(CASE WHEN event_type = 'unstake' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS rune_withdrawn,
	SUM(CASE WHEN event_type = 'add' THEN rune_amount ELSE 0 END) AS rune_added,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(units) AS units_changes,
    SUM(CASE WHEN event_type = 'rewards' THEN rune_amount ELSE 0 END) AS reward,
	SUM(CASE WHEN event_type = 'gas' THEN -asset_amount ELSE 0 END) AS gas_used,
	SUM(CASE WHEN event_type = 'gas' THEN rune_amount ELSE 0 END) AS gas_replenished,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY pool, time_bucket('1 day', time);

CREATE VIEW total_changes_5_min WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT time_bucket('5 minute', time) AS time,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN rune_amount > 0 AND event_type = 'swap' THEN rune_amount ELSE 0 END) AS buy_volume,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN rune_amount < 0 AND event_type = 'swap' THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(CASE WHEN event_type = 'rewards' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS total_reward,
    SUM(CASE WHEN event_type = 'rewards' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS total_deficit,
    COUNT(CASE WHEN event_type = 'add' THEN 1 ELSE NULL END) AS add_count,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY time_bucket('5 minute', time);

CREATE VIEW total_changes_hourly WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT time_bucket('1 hour', time) AS time,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN rune_amount > 0 AND event_type = 'swap' THEN rune_amount ELSE 0 END) AS buy_volume,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN rune_amount < 0 AND event_type = 'swap' THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(CASE WHEN event_type = 'rewards' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS total_reward,
    SUM(CASE WHEN event_type = 'rewards' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS total_deficit,
    COUNT(CASE WHEN event_type = 'add' THEN 1 ELSE NULL END) AS add_count,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY time_bucket('1 hour', time);

CREATE VIEW total_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT time_bucket('1 day', time) AS time,
    COUNT(CASE WHEN event_type = 'swap' AND asset_amount < 0 THEN 1 ELSE NULL END) AS buy_count,
    SUM(CASE WHEN rune_amount > 0 AND event_type = 'swap' THEN rune_amount ELSE 0 END) AS buy_volume,
    COUNT(CASE WHEN event_type = 'swap' AND rune_amount < 0 THEN 1 ELSE NULL END) AS sell_count,
    SUM(CASE WHEN rune_amount < 0 AND event_type = 'swap' THEN -rune_amount ELSE 0 END) AS sell_volume,
    SUM(CASE WHEN event_type = 'rewards' AND rune_amount > 0 THEN rune_amount ELSE 0 END) AS total_reward,
    SUM(CASE WHEN event_type = 'rewards' AND rune_amount < 0 THEN -rune_amount ELSE 0 END) AS total_deficit,
    COUNT(CASE WHEN event_type = 'add' THEN 1 ELSE NULL END) AS add_count,
    COUNT(CASE WHEN units > 0 THEN 1 ELSE NULL END) AS stake_count,
    COUNT(CASE WHEN units < 0 THEN 1 ELSE NULL END) AS withdraw_count
FROM pools_history
GROUP BY time_bucket('1 day', time);

CREATE VIEW stats_changes_5_min WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT time_bucket('5 minute', time) AS time,
    MIN(height) AS start_height,
    MAX(height) AS end_height,
    last(total_rune_depth, height) AS total_rune_depth,
    last(enabled_pools, height) AS enabled_pools,
    last(bootstrapped_pools, height) AS bootstrapped_pools,
    last(suspended_pools, height) AS suspended_pools
FROM stats_history
GROUP BY time_bucket('5 minute', time);

CREATE VIEW stats_changes_hourly WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT time_bucket('1 hour', time) AS time,
    MIN(height) AS start_height,
    MAX(height) AS end_height,
    last(total_rune_depth, height) AS total_rune_depth,
    last(enabled_pools, height) AS enabled_pools,
    last(bootstrapped_pools, height) AS bootstrapped_pools,
    last(suspended_pools, height) AS suspended_pools
FROM stats_history
GROUP BY time_bucket('1 hour', time);

CREATE VIEW stats_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_lag = "0", timescaledb.refresh_interval = '10 min')
AS
SELECT time_bucket('1 day', time) AS time,
    MIN(height) AS start_height,
    MAX(height) AS end_height,
    last(total_rune_depth, height) AS total_rune_depth,
    last(enabled_pools, height) AS enabled_pools,
    last(bootstrapped_pools, height) AS bootstrapped_pools,
    last(suspended_pools, height) AS suspended_pools
FROM stats_history
GROUP BY time_bucket('1 day', time);

-- +migrate Down

DROP VIEW pool_changes_5_min CASCADE;
DROP VIEW pool_changes_hourly CASCADE;
DROP VIEW pool_changes_daily CASCADE;
DROP VIEW total_changes_5_min CASCADE;
DROP VIEW total_changes_hourly CASCADE;
DROP VIEW total_changes_daily CASCADE;
DROP VIEW stats_changes_5_min CASCADE;
DROP VIEW stats_changes_hourly CASCADE;
DROP VIEW stats_changes_daily CASCADE;