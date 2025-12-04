package interfaces

import "github.com/rnegic/synchronous/pkg/telegramapi"

type TelegramAPIService interface {
	GetBotInfo() (*telegramapi.BotInfo, error)
	GetProfileByToken(accessToken string) (*telegramapi.BotInfo, error)
	SendMessage(chatID int64, text string) error
	SendMessageToUser(userID int64, message *telegramapi.SendMessageRequest) (*telegramapi.SendMessageResponse, error)
	GetChat(chatID int64) (*telegramapi.Chat, error)
	GetChatByLink(chatLink string) (*telegramapi.Chat, error)
	GetUserInfo(userID int64) (*telegramapi.TelegramUser, error)
	GetMessages(chatID int64, from, to, count *int64, messageIDs []string) ([]telegramapi.Message, error)
	AddMembers(chatID int64, userIDs []int64) error
	GetChatMembers(chatID int64, marker *int64, count *int, userIDs []int64) (*telegramapi.ChatMembersResponse, error)
	EditChat(chatID int64, title *string, icon interface{}) (*telegramapi.Chat, error)
	DeleteChat(chatID int64) error
	RemoveMember(chatID int64, userID int64) error
}
