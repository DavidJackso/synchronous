package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	AvatarURL *string        `gorm:"type:text" json:"avatarUrl"`
	MaxUserID int64          `gorm:"uniqueIndex:idx_max_user_id;not null" json:"maxUserId"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Stats *UserStats `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"stats,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type UserStats struct {
	UserID          string     `gorm:"type:varchar(36);primaryKey" json:"userId"`
	TotalSessions   int        `gorm:"not null;default:0" json:"totalSessions"`
	TotalFocusTime  int        `gorm:"not null;default:0" json:"totalFocusTime"` // в минутах
	CurrentStreak   int        `gorm:"not null;default:0" json:"currentStreak"`  // в днях
	LastSessionDate *time.Time `gorm:"type:date" json:"lastSessionDate,omitempty"`
	UpdatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (UserStats) TableName() string {
	return "user_stats"
}
