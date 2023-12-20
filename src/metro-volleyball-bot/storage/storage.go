package storage

type Storage interface {
	// GetLeaderboard will return the current leaderboard.
	ReadLeaderboard() (string, error)
	// SetLeaderboard will set the current leaderboard.
	UpdateLeaderboard(leaderboard string) error
	// Create Leaderboard will create a new leaderboard.
	CreateLeaderboard(leaderboard string) error
}

type LeaderBoardRecord struct {
}

type MemoryStorage struct {
	leaderboard string
}
