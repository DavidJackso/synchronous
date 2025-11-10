package service

import (
	"fmt"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type LeaderboardService struct {
	leaderboardRepo interfaces.LeaderboardRepository
	sessionRepo     interfaces.SessionRepository
	userRepo        interfaces.UserRepository
}

func NewLeaderboardService(
	leaderboardRepo interfaces.LeaderboardRepository,
	sessionRepo interfaces.SessionRepository,
	userRepo interfaces.UserRepository,
) interfaces.LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: leaderboardRepo,
		sessionRepo:     sessionRepo,
		userRepo:        userRepo,
	}
}

// GetSessionLeaderboard возвращает лидерборд для сессии
func (s *LeaderboardService) GetSessionLeaderboard(sessionID string, userID string) ([]*entity.LeaderboardEntry, error) {
	// Проверяем доступ к сессии
	_, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Получаем записи лидерборда из репозитория
	entries, err := s.leaderboardRepo.GetSessionLeaderboard(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	// Заполняем информацию о пользователях
	result := make([]*entity.LeaderboardEntry, 0, len(entries))
	for _, entry := range entries {
		user, err := s.userRepo.GetByID(entry.UserID)
		if err != nil {
			// Пропускаем пользователей, которых нет в БД
			continue
		}

		entry.UserName = user.Name
		entry.AvatarURL = user.AvatarURL
		result = append(result, entry)
	}

	return result, nil
}

// GetGlobalLeaderboard возвращает глобальный лидерборд
func (s *LeaderboardService) GetGlobalLeaderboard(userID string, period entity.LeaderboardPeriod, limit int) ([]*entity.LeaderboardEntry, error) {
	// Получаем записи лидерборда из репозитория
	entries, err := s.leaderboardRepo.GetGlobalLeaderboard(period, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get global leaderboard: %w", err)
	}

	// Заполняем информацию о пользователях
	result := make([]*entity.LeaderboardEntry, 0, len(entries))
	for _, entry := range entries {
		user, err := s.userRepo.GetByID(entry.UserID)
		if err != nil {
			// Пропускаем пользователей, которых нет в БД
			continue
		}

		entry.UserName = user.Name
		entry.AvatarURL = user.AvatarURL
		result = append(result, entry)
	}

	return result, nil
}
