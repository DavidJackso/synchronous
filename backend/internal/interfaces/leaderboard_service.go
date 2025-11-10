package interfaces

import (
	"github.com/rnegic/synchronous/internal/entity"
)

type LeaderboardService interface {
	GetSessionLeaderboard(sessionID string, userID string) ([]*entity.LeaderboardEntry, error)
	GetGlobalLeaderboard(userID string, period entity.LeaderboardPeriod, limit int) ([]*entity.LeaderboardEntry, error)
}
