package telegramapi

import (
	"encoding/json"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Update представляет обновление из Telegram Bot API (webhook)
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

// ParseUpdate парсит обновление из Telegram Bot API
func ParseUpdate(data []byte) (interface{}, error) {
	var update tgbotapi.Update
	if err := json.Unmarshal(data, &update); err != nil {
		return nil, fmt.Errorf("failed to parse update: %w", err)
	}

	switch {
	case update.Message != nil:
		if len(update.Message.NewChatMembers) > 0 {
			// Проверяем, был ли добавлен бот
			for _, member := range update.Message.NewChatMembers {
				if member.IsBot {
					return &BotAddedToChatUpdate{
						UpdateType: "bot_added",
						Timestamp:  int64(update.Message.Date),
						ChatID:     update.Message.Chat.ID,
						User: struct {
							UserID    int64  `json:"user_id"`
							FirstName string `json:"first_name"`
							Username  string `json:"username,omitempty"`
						}{
							UserID:    int64(update.Message.From.ID),
							FirstName: update.Message.From.FirstName,
							Username:  update.Message.From.UserName,
						},
						IsChannel: update.Message.Chat.Type == "channel",
					}, nil
				}
			}
		}

		// Обработка команды /start с параметром
		if update.Message.Text != "" && len(update.Message.Text) > 6 && update.Message.Text[:6] == "/start" {
			var startPayload string
			if len(update.Message.Text) > 7 {
				startPayload = update.Message.Text[7:]
			}

			if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
				chat := convertChatFromUpdate(update.Message.Chat)
				return &MessageChatCreatedUpdate{
					UpdateType:   "message_chat_created",
					Timestamp:    int64(update.Message.Date),
					Chat:         *chat,
					MessageID:    "",
					StartPayload: startPayload,
				}, nil
			}
		}

		// Обычное сообщение
		return &MessageCreatedUpdate{
			UpdateType: "message_created",
			Timestamp:  int64(update.Message.Date),
			Message:    convertMessage(*update.Message),
		}, nil

	case update.CallbackQuery != nil:
		return &MessageCallbackUpdate{
			UpdateType: "message_callback",
			Timestamp:  int64(update.CallbackQuery.Message.Date),
			Callback: Callback{
				Timestamp:  int64(update.CallbackQuery.Message.Date),
				CallbackID: update.CallbackQuery.ID,
				Payload:    update.CallbackQuery.Data,
				User: struct {
					UserID    int64  `json:"user_id"`
					FirstName string `json:"first_name"`
					Username  string `json:"username,omitempty"`
				}{
					UserID:    int64(update.CallbackQuery.From.ID),
					FirstName: update.CallbackQuery.From.FirstName,
					Username:  update.CallbackQuery.From.UserName,
				},
			},
			Message: func() *Message {
				msg := convertMessage(*update.CallbackQuery.Message)
				return &msg
			}(),
		}, nil

	default:
		return &Update{
			UpdateType: "unknown",
			Timestamp:  time.Now().Unix(),
		}, nil
	}
}

// ToTime конвертирует Unix timestamp в time.Time
func (u *Update) ToTime() time.Time {
	return time.Unix(u.Timestamp, 0)
}

func convertChatFromUpdate(chat *tgbotapi.Chat) *Chat {
	chatType := string(chat.Type)
	status := "active"

	// В Telegram Bot API количество участников не всегда доступно
	participantsCount := 0

	return &Chat{
		ChatID:            chat.ID,
		Type:              chatType,
		Status:            status,
		Title:             chat.Title,
		LastEventTime:     0,
		ParticipantsCount: participantsCount,
	}
}
