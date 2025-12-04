package entity

// TelegramChatInfo информация о чате в Telegram API
type TelegramChatInfo struct {
	ChatID            int64  `json:"chatId"`
	ChatLink          string `json:"chatLink"`
	Title             string `json:"title"`
	ParticipantsCount int    `json:"participantsCount"`
}
