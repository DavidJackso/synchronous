package service

import (
	"fmt"
	"time"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type MessageService struct {
	sessionService     interfaces.SessionService
	telegramAPIService interfaces.TelegramAPIService
	userRepo           interfaces.UserRepository
	messageRepo        interfaces.MessageRepository
}

func NewMessageService(
	sessionService interfaces.SessionService,
	telegramAPIService interfaces.TelegramAPIService,
	userRepo interfaces.UserRepository,
	messageRepo interfaces.MessageRepository,
) interfaces.MessageService {
	return &MessageService{
		sessionService:     sessionService,
		telegramAPIService: telegramAPIService,
		userRepo:           userRepo,
		messageRepo:        messageRepo,
	}
}

// GetMessages получает сообщения из Telegram API для сессии
func (s *MessageService) GetMessages(sessionID string, userID string, before *time.Time, limit int) ([]*entity.Message, error) {
	// Проверяем доступ к сессии
	session, err := s.sessionService.GetSession(sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Проверяем, что чат создан
	if session.TelegramChatID == nil {
		return nil, fmt.Errorf("chat not created for this session")
	}

	// Преобразуем before в Unix timestamp в миллисекундах для Telegram API
	var to *int64
	if before != nil {
		timestamp := before.Unix() * 1000
		to = &timestamp
	}

	count := int64(limit)

	// Получаем сообщения из Telegram API
	telegramMessages, err := s.telegramAPIService.GetMessages(*session.TelegramChatID, nil, to, &count, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages from Telegram API: %w", err)
	}

	// Преобразуем сообщения Telegram API в entity.Message
	messages := make([]*entity.Message, 0, len(telegramMessages))
	for _, telegramMsg := range telegramMessages {
		// Получаем информацию о пользователе
		user, err := s.userRepo.GetByTelegramUserID(telegramMsg.Sender.UserID)
		if err != nil {
			// Пропускаем сообщения от пользователей, которых нет в нашей БД
			continue
		}

		msg := &entity.Message{
			ID:        telegramMsg.Body.Mid,
			UserID:    user.ID,
			UserName:  user.Name,
			AvatarURL: user.AvatarURL,
			Text:      telegramMsg.Body.Text,
			SessionID: sessionID,
			CreatedAt: time.Unix(telegramMsg.Timestamp/1000, (telegramMsg.Timestamp%1000)*1000000),
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// SendMessage отправляет сообщение через Telegram API
func (s *MessageService) SendMessage(sessionID string, userID string, text string) (*entity.Message, error) {
	// Проверяем доступ к сессии
	session, err := s.sessionService.GetSession(sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Проверяем, что чат создан
	if session.TelegramChatID == nil {
		return nil, fmt.Errorf("chat not created for this session")
	}

	// Получаем информацию о пользователе
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Отправляем сообщение через Telegram API
	// TelegramUserID обязательное поле (not null), поэтому проверка не нужна
	// Отправляем сообщение от имени бота в чат
	err = s.telegramAPIService.SendMessage(*session.TelegramChatID, text)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to Telegram API: %w", err)
	}

	// Создаем объект сообщения для ответа
	// В реальности Telegram API должен вернуть информацию о созданном сообщении
	// Пока создаем простой объект
	msg := &entity.Message{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()), // Временный ID
		UserID:    user.ID,
		UserName:  user.Name,
		AvatarURL: user.AvatarURL,
		Text:      text,
		SessionID: sessionID,
		CreatedAt: time.Now(),
	}

	return msg, nil
}

// GetChatInfo возвращает информацию о чате Telegram для сессии
func (s *MessageService) GetChatInfo(sessionID string, userID string) (*entity.TelegramChatInfo, error) {
	// Проверяем доступ к сессии
	session, err := s.sessionService.GetSession(sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Проверяем, что чат создан
	if session.TelegramChatID == nil {
		return nil, fmt.Errorf("chat not created for this session")
	}

	// Получаем информацию о чате из Telegram API
	chat, err := s.telegramAPIService.GetChat(*session.TelegramChatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat info from Telegram API: %w", err)
	}

	chatInfo := &entity.TelegramChatInfo{
		ChatID:            chat.ChatID,
		ChatLink:          "", // В реальности нужно получить из Telegram API или сформировать
		Title:             chat.Title,
		ParticipantsCount: chat.ParticipantsCount,
	}

	// Если есть ссылка в сессии, используем её
	if session.TelegramChatLink != nil {
		chatInfo.ChatLink = *session.TelegramChatLink
	}

	return chatInfo, nil
}
