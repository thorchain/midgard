package repository

// LatestState represents latest height and eventID inserted database.
type LatestState struct {
	Height  int64 `db:"height"`
	EventID int64 `db:"event_id"`
}
