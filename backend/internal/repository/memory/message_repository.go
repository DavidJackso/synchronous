package memory

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type MessageRepository struct {
	messages map[string]*entity.Message
	mu       sync.RWMutex
}

func NewMessageRepository() interfaces.MessageRepository {
	return &MessageRepository{
		messages: make(map[string]*entity.Message),
	}
}

func (r *MessageRepository) Create(message *entity.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.messages[message.ID]; exists {
		return fmt.Errorf("message with ID %s already exists", message.ID)
	}

	r.messages[message.ID] = message
	return nil
}

func (r *MessageRepository) GetBySessionID(sessionID string, before *time.Time, limit int) ([]*entity.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var messages []*entity.Message
	for _, message := range r.messages {
		if message.SessionID == sessionID {
			if before == nil || message.CreatedAt.Before(*before) {
				messages = append(messages, message)
			}
		}
	}

	// Сортируем по времени создания (от новых к старым)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})

	// Ограничиваем количество
	if limit > 0 && len(messages) > limit {
		messages = messages[:limit]
	}

	return messages, nil
}

func (r *MessageRepository) GetByID(id string) (*entity.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	message, exists := r.messages[id]
	if !exists {
		return nil, fmt.Errorf("message with ID %s not found", id)
	}

	return message, nil
}
