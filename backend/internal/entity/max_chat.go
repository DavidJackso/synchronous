package entity

// MaxChatInfo информация о чате в Max API
type MaxChatInfo struct {
	ChatID            int64  `json:"chatId"`
	ChatLink          string `json:"chatLink"`
	Title             string `json:"title"`
	ParticipantsCount int    `json:"participantsCount"`
}
