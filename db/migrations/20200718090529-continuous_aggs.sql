
-- +migrate Up

CREATE VIEW pool_event_changes_daily WITH
(timescaledb.continuous, timescaledb.refresh_interval = '5s')
AS
SELECT pool, event_type, time_bucket('1 day', time) AS time,
    SUM(CASE WHEN asset_amount > 0 THEN asset_amount ELSE 0 END) AS pos_asset_changes,
    SUM(CASE WHEN asset_amount < 0 THEN asset_amount ELSE 0 END) AS neg_asset_changes,
    SUM(CASE WHEN rune_amount > 0 THEN rune_amount ELSE 0 END) AS pos_rune_changes,
    SUM(CASE WHEN rune_amount < 0 THEN rune_amount ELSE 0 END) AS neg_rune_changes,
    SUM(units) AS units_changes
FROM pools_history
GROUP BY pool, event_type, time_bucket('1 day', time);

-- +migrate Down

DROP VIEW pool_event_changes_daily CASCADE;