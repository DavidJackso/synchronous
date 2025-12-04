package entity

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID                string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	SessionID         string         `gorm:"type:varchar(36);not null;index:idx_session_id" json:"sessionId"`
	UserID            string         `gorm:"type:varchar(36);not null;index:idx_user_id" json:"userId"`
	UserName          string         `gorm:"type:varchar(255);not null" json:"userName"`
	AvatarURL         *string        `gorm:"type:text" json:"avatarUrl"`
	Text              string         `gorm:"type:text;not null" json:"text"`
	TelegramMessageID *string        `gorm:"type:varchar(255);index:idx_telegram_message_id" json:"telegramMessageId,omitempty"` // ID сообщения в Telegram API
	CreatedAt         time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_created_at" json:"createdAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Session *Session `gorm:"foreignKey:SessionID" json:"session,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
