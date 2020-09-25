
-- +migrate Up

CREATE VIEW pool_changes_5_min WITH (timescaledb.continuous) AS
SELECT 
    pool,
    time_bucket('5 min', time) AS time,
    first(height, id) AS start_height,
    last(height, id) AS end_height,
    SUM(asset_amount) AS asset_changes,
    SUM(
        CASE
            WHEN event_type = 'stake'::event_type THEN asset_amount
            ELSE 0
        END
    ) AS asset_staked,
    SUM(
        CASE
            WHEN event_type = 'unstake'::event_type
            AND asset_amount < 0 THEN -asset_amount
            ELSE 0
        END
    ) AS asset_withdrawn,
    COUNT(
        CASE
            WHEN event_type = 'swap'::event_type
            AND asset_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS buy_count,
    SUM(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount > 0 THEN rune_amount
            ELSE 0
        END
    ) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    SUM(
        CASE
            WHEN event_type = 'stake'::event_type THEN rune_amount
            ELSE 0
        END
    ) AS rune_staked,
    SUM(
        CASE
            WHEN event_type = 'unstake'::event_type
            AND rune_amount < 0 THEN -rune_amount
            ELSE 0
        END
    ) AS rune_withdrawn,
    COUNT(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS sell_count,
    SUM(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount < 0 THEN -rune_amount
            ELSE 0
        END
    ) AS sell_volume,
    SUM((meta->>'units')::BIGINT) AS units_changes,
    COUNT(
        CASE
            WHEN type = 'stake'::event_type
            AND (meta->>'units')::BIGINT > 0 THEN 1
            ELSE NULL
        END
    ) AS stake_count,
    COUNT(
        CASE
            WHEN type = 'unstake'::event_type
            THEN 1
            ELSE NULL
        END
    ) AS withdraw_count
FROM events
WHERE event_status = 'success'::event_status
GROUP BY 
    pool,
    time_bucket('5 min', time);

CREATE VIEW pool_changes_1_hour WITH (timescaledb.continuous) AS
SELECT 
    pool,
    time_bucket('1 hour', time) AS time,
    first(height, id) AS start_height,
    last(height, id) AS end_height,
    SUM(asset_amount) AS asset_changes,
    SUM(
        CASE
            WHEN event_type = 'stake'::event_type THEN asset_amount
            ELSE 0
        END
    ) AS asset_staked,
    SUM(
        CASE
            WHEN event_type = 'unstake'::event_type
            AND asset_amount < 0 THEN -asset_amount
            ELSE 0
        END
    ) AS asset_withdrawn,
    COUNT(
        CASE
            WHEN event_type = 'swap'::event_type
            AND asset_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS buy_count,
    SUM(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount > 0 THEN rune_amount
            ELSE 0
        END
    ) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    SUM(
        CASE
            WHEN event_type = 'stake'::event_type THEN rune_amount
            ELSE 0
        END
    ) AS rune_staked,
    SUM(
        CASE
            WHEN event_type = 'unstake'::event_type
            AND rune_amount < 0 THEN -rune_amount
            ELSE 0
        END
    ) AS rune_withdrawn,
    COUNT(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS sell_count,
    SUM(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount < 0 THEN -rune_amount
            ELSE 0
        END
    ) AS sell_volume,
    SUM((meta->>'units')::BIGINT) AS units_changes,
    COUNT(
        CASE
            WHEN type = 'stake'::event_type
            AND (meta->>'units')::BIGINT > 0 THEN 1
            ELSE NULL
        END
    ) AS stake_count,
    COUNT(
        CASE
            WHEN type = 'unstake'::event_type
            THEN 1
            ELSE NULL
        END
    ) AS withdraw_count
FROM events
WHERE event_status = 'success'::event_status
GROUP BY 
    pool,
    time_bucket('1 hour', time);

CREATE VIEW pool_changes_1_day WITH (timescaledb.continuous) AS
SELECT 
    pool,
    time_bucket('1 day', time) AS time,
    first(height, id) AS start_height,
    last(height, id) AS end_height,
    SUM(asset_amount) AS asset_changes,
    SUM(
        CASE
            WHEN event_type = 'stake'::event_type THEN asset_amount
            ELSE 0
        END
    ) AS asset_staked,
    SUM(
        CASE
            WHEN event_type = 'unstake'::event_type
            AND asset_amount < 0 THEN -asset_amount
            ELSE 0
        END
    ) AS asset_withdrawn,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND asset_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS buy_count,
    SUM(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount > 0 THEN rune_amount
            ELSE 0
        END
    ) AS buy_volume,
    SUM(rune_amount) AS rune_changes,
    SUM(
        CASE
            WHEN event_type = 'stake'::event_type THEN rune_amount
            ELSE 0
        END
    ) AS rune_staked,
    SUM(
        CASE
            WHEN event_type = 'unstake'::event_type
            AND rune_amount < 0 THEN -rune_amount
            ELSE 0
        END
    ) AS rune_withdrawn,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND rune_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS sell_count,
    SUM(
        CASE
            WHEN event_type = 'swap'::event_type
            AND rune_amount < 0 THEN -rune_amount
            ELSE 0
        END
    ) AS sell_volume,
    SUM((meta->>'units')::BIGINT) AS units_changes,
    COUNT(
        CASE
            WHEN type = 'stake'::event_type
            AND (meta->>'units')::BIGINT > 0 THEN 1
            ELSE NULL
        END
    ) AS stake_count,
    COUNT(
        CASE
            WHEN type = 'unstake'::event_type
            THEN 1
            ELSE NULL
        END
    ) AS withdraw_count
FROM events
WHERE event_status = 'success'::event_status
GROUP BY 
    pool,
    time_bucket('1 day', time);

CREATE VIEW stats_changes_5_min WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('5 minute', time) AS time,
    first(height, id) AS start_height,
    last(height, id) AS end_height,
    SUM(rune_amount) AS rune_changes,
    COUNT(tx_hash) AS txs_count,
    SUM(
        CASE
            WHEN rune_amount > 0
            AND event_type = 'swap'::event_type THEN rune_amount
            ELSE 0
        END
    ) AS buy_volume,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND rune_amount > 0 THEN 1
            ELSE NULL
        END
    ) AS buy_count,
    SUM(
        CASE
            WHEN rune_amount < 0
            AND event_type = 'swap'::event_type THEN -rune_amount
            ELSE 0
        END
    ) AS sell_volume,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND rune_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS sell_count,
    COUNT(
        CASE
            WHEN type = 'stake'::event_type
            AND (meta->>'units')::BIGINT > 0 THEN 1
            ELSE NULL
        END
    ) AS stake_count,
    COUNT(
        CASE
            WHEN type = 'unstake'::event_type
            THEN 1
            ELSE NULL
        END
    ) AS withdraw_count
FROM events
WHERE event_status = 'success'::event_status
GROUP BY time_bucket('5 minute', time);


CREATE VIEW stats_changes_1_hour WITH (timescaledb.continuous) AS
SELECT time_bucket('1 hour', time) AS time,
    first(height, id) AS start_height,
    last(height, id) AS end_height,
    SUM(rune_amount) AS rune_changes,
    COUNT(tx_hash) AS txs_count,
    SUM(
        CASE
            WHEN rune_amount > 0
            AND event_type = 'swap'::event_type THEN rune_amount
            ELSE 0
        END
    ) AS buy_volume,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND rune_amount > 0 THEN 1
            ELSE NULL
        END
    ) AS buy_count,
    SUM(
        CASE
            WHEN rune_amount < 0
            AND event_type = 'swap'::event_type THEN -rune_amount
            ELSE 0
        END
    ) AS sell_volume,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND rune_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS sell_count,
    COUNT(
        CASE
            WHEN type = 'stake'::event_type
            AND (meta->>'units')::BIGINT > 0 THEN 1
            ELSE NULL
        END
    ) AS stake_count,
    COUNT(
        CASE
            WHEN type = 'unstake'::event_type
            THEN 1
            ELSE NULL
        END
    ) AS withdraw_count
FROM events
WHERE event_status = 'success'::event_status
GROUP BY time_bucket('1 hour', time);

CREATE VIEW stats_changes_1_day WITH (timescaledb.continuous) AS
SELECT time_bucket('1 day', time) AS time,
    first(height, id) AS start_height,
    last(height, id) AS end_height,
    SUM(rune_amount) AS rune_changes,
    COUNT(tx_hash) AS txs_count,
    SUM(
        CASE
            WHEN rune_amount > 0
            AND event_type = 'swap'::event_type THEN rune_amount
            ELSE 0
        END
    ) AS buy_volume,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND rune_amount > 0 THEN 1
            ELSE NULL
        END
    ) AS buy_count,
    SUM(
        CASE
            WHEN rune_amount < 0
            AND event_type = 'swap'::event_type THEN -rune_amount
            ELSE 0
        END
    ) AS sell_volume,
    COUNT(
        CASE
            WHEN type = 'swap'::event_type
            AND rune_amount < 0 THEN 1
            ELSE NULL
        END
    ) AS sell_count,
    COUNT(
        CASE
            WHEN type = 'stake'::event_type
            AND (meta->>'units')::BIGINT > 0 THEN 1
            ELSE NULL
        END
    ) AS stake_count,
    COUNT(
        CASE
            WHEN type = 'unstake'::event_type
            THEN 1
            ELSE NULL
        END
    ) AS withdraw_count
FROM events
WHERE event_status = 'success'::event_status
GROUP BY time_bucket('1 day', time);
-- +migrate Down
DROP VIEW pool_changes_5_min CASCADE;
DROP VIEW pool_changes_1_hour CASCADE;
DROP VIEW pool_changes_1_day CASCADE;
DROP VIEW stats_changes_5_min CASCADE;
DROP VIEW stats_changes_1_hour CASCADE;
DROP VIEW stats_changes_1_day CASCADE;