package interfaces

import (
	"time"

	"github.com/rnegic/synchronous/internal/entity"
)

type MessageService interface {
	// GetMessages получает сообщения из Telegram API (не из нашей БД)
	// Сообщения хранятся в Telegram, мы их не дублируем
	GetMessages(sessionID string, userID string, before *time.Time, limit int) ([]*entity.Message, error)

	// SendMessage отправляет сообщение через Telegram API (не сохраняет в БД)
	// Сообщение отправляется в чат Telegram и хранится там
	SendMessage(sessionID string, userID string, text string) (*entity.Message, error)

	// GetChatInfo возвращает информацию о чате Telegram для сессии
	GetChatInfo(sessionID string, userID string) (*entity.TelegramChatInfo, error)
}
