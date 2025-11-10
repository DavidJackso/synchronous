package memory

import (
	"sort"
	"sync"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type LeaderboardRepository struct {
	sessionScores map[string]map[string]int // sessionID -> userID -> score
	globalScores  map[string]int            // userID -> totalScore
	mu            sync.RWMutex
}

func NewLeaderboardRepository() interfaces.LeaderboardRepository {
	return &LeaderboardRepository{
		sessionScores: make(map[string]map[string]int),
		globalScores:  make(map[string]int),
	}
}

func (r *LeaderboardRepository) GetSessionLeaderboard(sessionID string) ([]*entity.LeaderboardEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	scores, exists := r.sessionScores[sessionID]
	if !exists {
		return []*entity.LeaderboardEntry{}, nil
	}

	type userScore struct {
		userID string
		score  int
	}

	var userScores []userScore
	for userID, score := range scores {
		userScores = append(userScores, userScore{userID: userID, score: score})
	}

	// Сортируем по убыванию score
	sort.Slice(userScores, func(i, j int) bool {
		return userScores[i].score > userScores[j].score
	})

	// Формируем результат
	var entries []*entity.LeaderboardEntry
	for rank, us := range userScores {
		entries = append(entries, &entity.LeaderboardEntry{
			Rank:   rank + 1,
			UserID: us.userID,
			Score:  us.score,
		})
	}

	return entries, nil
}

func (r *LeaderboardRepository) GetGlobalLeaderboard(period entity.LeaderboardPeriod, limit int) ([]*entity.LeaderboardEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type userScore struct {
		userID string
		score  int
	}

	var userScores []userScore
	for userID, score := range r.globalScores {
		userScores = append(userScores, userScore{userID: userID, score: score})
	}

	// Сортируем по убыванию score
	sort.Slice(userScores, func(i, j int) bool {
		return userScores[i].score > userScores[j].score
	})

	// Ограничиваем количество
	if limit > 0 && len(userScores) > limit {
		userScores = userScores[:limit]
	}

	// Формируем результат
	var entries []*entity.LeaderboardEntry
	for rank, us := range userScores {
		entries = append(entries, &entity.LeaderboardEntry{
			Rank:   rank + 1,
			UserID: us.userID,
			Score:  us.score,
		})
	}

	return entries, nil
}

func (r *LeaderboardRepository) UpdateUserScore(userID string, sessionID string, score int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Обновляем счет для сессии
	if r.sessionScores[sessionID] == nil {
		r.sessionScores[sessionID] = make(map[string]int)
	}
	r.sessionScores[sessionID][userID] = score

	// Обновляем глобальный счет
	r.globalScores[userID] += score

	return nil
}
