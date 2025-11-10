package interfaces

import "github.com/rnegic/synchronous/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByMaxUserID(maxUserID int64) (*entity.User, error)
	Update(user *entity.User) error
	UpdateStats(userID string, stats *entity.UserStats) error
	GetStats(userID string) (*entity.UserStats, error)
}
