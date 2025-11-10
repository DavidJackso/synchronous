package memory

import (
	"fmt"
	"sync"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type UserRepository struct {
	users map[string]*entity.User
	stats map[string]*entity.UserStats
	maxID map[int64]string // маппинг maxUserID -> userID
	mu    sync.RWMutex
}

func NewUserRepository() interfaces.UserRepository {
	return &UserRepository{
		users: make(map[string]*entity.User),
		stats: make(map[string]*entity.UserStats),
		maxID: make(map[int64]string),
	}
}

func (r *UserRepository) Create(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}

	r.users[user.ID] = user
	r.maxID[user.MaxUserID] = user.ID
	r.stats[user.ID] = &entity.UserStats{}

	return nil
}

func (r *UserRepository) GetByID(id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}

	return user, nil
}

func (r *UserRepository) GetByMaxUserID(maxUserID int64) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.maxID[maxUserID]
	if !exists {
		return nil, fmt.Errorf("user with maxUserID %d not found", maxUserID)
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", userID)
	}

	return user, nil
}

func (r *UserRepository) Update(user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user with ID %s not found", user.ID)
	}

	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) UpdateStats(userID string, stats *entity.UserStats) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[userID]; !exists {
		return fmt.Errorf("user with ID %s not found", userID)
	}

	r.stats[userID] = stats
	return nil
}

func (r *UserRepository) GetStats(userID string) (*entity.UserStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats, exists := r.stats[userID]
	if !exists {
		return &entity.UserStats{}, nil
	}

	return stats, nil
}
