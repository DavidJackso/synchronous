package gorm

import (
	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByMaxUserID(maxUserID int64) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("max_user_id = ?", maxUserID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateStats(userID string, stats *entity.UserStats) error {
	stats.UserID = userID
	return r.db.Where("user_id = ?", userID).Save(stats).Error
}

func (r *userRepository) GetStats(userID string) (*entity.UserStats, error) {
	var stats entity.UserStats
	err := r.db.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Создаем пустую статистику, если её нет
			stats = entity.UserStats{
				UserID: userID,
			}
			if err := r.db.Create(&stats).Error; err != nil {
				return nil, err
			}
			return &stats, nil
		}
		return nil, err
	}
	return &stats, nil
}
