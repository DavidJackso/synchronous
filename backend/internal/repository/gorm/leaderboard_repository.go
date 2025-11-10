package gorm

import (
	"time"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
	"gorm.io/gorm"
)

type leaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) interfaces.LeaderboardRepository {
	return &leaderboardRepository{db: db}
}

func (r *leaderboardRepository) GetSessionLeaderboard(sessionID string) ([]*entity.LeaderboardEntry, error) {
	var entries []*entity.LeaderboardEntry

	err := r.db.Table("session_participants").
		Select(`
			session_participants.user_id,
			session_participants.user_name,
			session_participants.avatar_url,
			COALESCE(COUNT(CASE WHEN tasks.completed = true THEN 1 END), 0) as tasks_completed,
			COALESCE(SUM(sessions.focus_duration), 0) as focus_time
		`).
		Joins("JOIN sessions ON session_participants.session_id = sessions.id").
		Joins("LEFT JOIN tasks ON tasks.session_id = sessions.id").
		Where("session_participants.session_id = ?", sessionID).
		Group("session_participants.user_id, session_participants.user_name, session_participants.avatar_url").
		Order("tasks_completed DESC, focus_time DESC").
		Scan(&entries).Error

	if err != nil {
		return nil, err
	}

	// Добавляем rank и score
	for i := range entries {
		entries[i].Rank = i + 1
		entries[i].Score = entries[i].TasksCompleted*10 + entries[i].FocusTime
	}

	return entries, nil
}

func (r *leaderboardRepository) GetGlobalLeaderboard(period entity.LeaderboardPeriod, limit int) ([]*entity.LeaderboardEntry, error) {
	var entries []*entity.LeaderboardEntry
	query := r.db.Table("users").
		Select(`
			users.id as user_id,
			users.name as user_name,
			users.avatar_url,
			COALESCE(user_stats.total_sessions, 0) as tasks_completed,
			COALESCE(user_stats.total_focus_time, 0) as focus_time
		`).
		Joins("LEFT JOIN user_stats ON users.id = user_stats.user_id")

	// Фильтр по периоду
	now := time.Now()
	switch period {
	case entity.LeaderboardPeriodDay:
		query = query.Where("user_stats.updated_at >= ?", now.AddDate(0, 0, -1))
	case entity.LeaderboardPeriodWeek:
		query = query.Where("user_stats.updated_at >= ?", now.AddDate(0, 0, -7))
	case entity.LeaderboardPeriodMonth:
		query = query.Where("user_stats.updated_at >= ?", now.AddDate(0, -1, 0))
		// LeaderboardPeriodAll - без фильтра
	}

	err := query.
		Order("tasks_completed DESC, focus_time DESC").
		Limit(limit).
		Scan(&entries).Error

	if err != nil {
		return nil, err
	}

	// Добавляем rank и score
	for i := range entries {
		entries[i].Rank = i + 1
		entries[i].Score = entries[i].TasksCompleted*10 + entries[i].FocusTime
	}

	return entries, nil
}

func (r *leaderboardRepository) UpdateUserScore(userID string, sessionID string, score int) error {
	// Обновляем статистику пользователя на основе завершенной сессии
	var session entity.Session
	if err := r.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return err
	}

	// Подсчитываем выполненные задачи
	var tasksCompleted int64
	r.db.Model(&entity.Task{}).
		Where("session_id = ? AND completed = ?", sessionID, true).
		Count(&tasksCompleted)

	// Обновляем или создаем статистику
	stats := entity.UserStats{
		UserID:         userID,
		TotalSessions:  1,
		TotalFocusTime: session.FocusDuration,
	}

	err := r.db.Where("user_id = ?", userID).FirstOrCreate(&stats).Error
	if err != nil {
		return err
	}

	// Обновляем статистику
	return r.db.Model(&stats).
		Updates(map[string]interface{}{
			"total_sessions":   gorm.Expr("total_sessions + 1"),
			"total_focus_time": gorm.Expr("total_focus_time + ?", session.FocusDuration),
			"updated_at":       time.Now(),
		}).Error
}
