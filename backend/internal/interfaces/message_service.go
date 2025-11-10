package interfaces

import (
	"time"

	"github.com/rnegic/synchronous/internal/entity"
)

type MessageService interface {
	// GetMessages получает сообщения из Max API (не из нашей БД)
	// Сообщения хранятся в Max, мы их не дублируем
	GetMessages(sessionID string, userID string, before *time.Time, limit int) ([]*entity.Message, error)

	// SendMessage отправляет сообщение через Max API (не сохраняет в БД)
	// Сообщение отправляется в чат Max и хранится там
	SendMessage(sessionID string, userID string, text string) (*entity.Message, error)

	// GetChatInfo возвращает информацию о чате Max для сессии
	GetChatInfo(sessionID string, userID string) (*entity.MaxChatInfo, error)
}
