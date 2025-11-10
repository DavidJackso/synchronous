package maxapi

import (
	"encoding/json"
	"fmt"
	"time"
)

// Update представляет обновление из Max API (webhook или long polling)
type Update struct {
	UpdateType string          `json:"update_type"`
	Timestamp  int64           `json:"timestamp"`
	RawData    json.RawMessage `json:"-"` // Для хранения полных данных
}

// MessageCreatedUpdate обновление о создании сообщения
type MessageCreatedUpdate struct {
	UpdateType string  `json:"update_type"`
	Timestamp  int64   `json:"timestamp"`
	Message    Message `json:"message"`
	UserLocale *string `json:"user_locale,omitempty"`
}

// MessageCallbackUpdate обновление о нажатии на кнопку
type MessageCallbackUpdate struct {
	UpdateType string   `json:"update_type"`
	Timestamp  int64    `json:"timestamp"`
	Callback   Callback `json:"callback"`
	Message    *Message `json:"message,omitempty"`
	UserLocale *string  `json:"user_locale,omitempty"`
}

// Callback представляет callback от кнопки
type Callback struct {
	Timestamp  int64  `json:"timestamp"`
	CallbackID string `json:"callback_id"`
	Payload    string `json:"payload,omitempty"`
	User       struct {
		UserID    int64  `json:"user_id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username,omitempty"`
	} `json:"user"`
}

// BotAddedToChatUpdate обновление о добавлении бота в чат
type BotAddedToChatUpdate struct {
	UpdateType string `json:"update_type"`
	Timestamp  int64  `json:"timestamp"`
	ChatID     int64  `json:"chat_id"`
	User       struct {
		UserID    int64  `json:"user_id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username,omitempty"`
	} `json:"user"`
	IsChannel bool `json:"is_channel"`
}

// MessageChatCreatedUpdate обновление о создании чата через кнопку
type MessageChatCreatedUpdate struct {
	UpdateType   string `json:"update_type"`
	Timestamp    int64  `json:"timestamp"`
	Chat         Chat   `json:"chat"`
	MessageID    string `json:"message_id"`
	StartPayload string `json:"start_payload,omitempty"`
}

// ParseUpdate парсит обновление из Max API
func ParseUpdate(data []byte) (interface{}, error) {
	var baseUpdate Update
	if err := json.Unmarshal(data, &baseUpdate); err != nil {
		return nil, fmt.Errorf("failed to parse update: %w", err)
	}

	baseUpdate.RawData = data

	switch baseUpdate.UpdateType {
	case "message_created":
		var update MessageCreatedUpdate
		if err := json.Unmarshal(data, &update); err != nil {
			return nil, fmt.Errorf("failed to parse message_created update: %w", err)
		}
		return &update, nil

	case "message_callback":
		var update MessageCallbackUpdate
		if err := json.Unmarshal(data, &update); err != nil {
			return nil, fmt.Errorf("failed to parse message_callback update: %w", err)
		}
		return &update, nil

	case "bot_added":
		var update BotAddedToChatUpdate
		if err := json.Unmarshal(data, &update); err != nil {
			return nil, fmt.Errorf("failed to parse bot_added update: %w", err)
		}
		return &update, nil

	case "message_chat_created":
		var update MessageChatCreatedUpdate
		if err := json.Unmarshal(data, &update); err != nil {
			return nil, fmt.Errorf("failed to parse message_chat_created update: %w", err)
		}
		return &update, nil

	default:
		// Возвращаем базовое обновление для неизвестных типов
		return &baseUpdate, nil
	}
}

// ToTime конвертирует Unix timestamp в time.Time
func (u *Update) ToTime() time.Time {
	return time.Unix(u.Timestamp/1000, (u.Timestamp%1000)*1000000)
}
