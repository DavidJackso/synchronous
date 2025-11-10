package gorm

import (
	"time"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) interfaces.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *entity.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetBySessionID(sessionID string, before *time.Time, limit int) ([]*entity.Message, error) {
	var messages []*entity.Message
	query := r.db.Where("session_id = ?", sessionID)

	if before != nil {
		query = query.Where("created_at < ?", *before)
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, err
	}

	// Разворачиваем порядок для правильной последовательности
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepository) GetByID(id string) (*entity.Message, error) {
	var message entity.Message
	err := r.db.Where("id = ?", id).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}
