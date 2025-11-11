package interfaces

import "github.com/rnegic/synchronous/pkg/maxapi"

type MaxAPIService interface {
	GetBotInfo() (*maxapi.BotInfo, error)
	GetProfileByToken(accessToken string) (*maxapi.BotInfo, error)
	SendMessage(chatID int64, text string) error
	SendMessageToUser(userID int64, message *maxapi.SendMessageRequest) (*maxapi.SendMessageResponse, error)
	GetChat(chatID int64) (*maxapi.Chat, error)
	GetChatByLink(chatLink string) (*maxapi.Chat, error)
	GetUserInfo(userID int64) (*maxapi.MaxUser, error)
	GetMessages(chatID int64, from, to, count *int64, messageIDs []string) ([]maxapi.Message, error)
	AddMembers(chatID int64, userIDs []int64) error
	GetChatMembers(chatID int64, marker *int64, count *int, userIDs []int64) (*maxapi.ChatMembersResponse, error)
	EditChat(chatID int64, title *string, icon interface{}) (*maxapi.Chat, error)
	DeleteChat(chatID int64) error
	RemoveMember(chatID int64, userID int64) error
}
