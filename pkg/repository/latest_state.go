package repository

// LatestState represents latest height and eventID inserted database.
type LatestState struct {
	Height  int64
	EventID int64
}
