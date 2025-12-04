package service

import (
	"fmt"
	"log"

	"github.com/rnegic/synchronous/internal/interfaces"
	"github.com/rnegic/synchronous/pkg/telegramapi"
)

type TelegramAPIService struct {
	client *telegramapi.Client
}

func NewTelegramAPIService(botToken string) interfaces.TelegramAPIService {
	client, err := telegramapi.NewClient(botToken)
	if err != nil {
		log.Fatalf("failed to initialize Telegram API client: %v", err)
	}

	return &TelegramAPIService{
		client: client,
	}
}

func (s *TelegramAPIService) GetBotInfo() (*telegramapi.BotInfo, error) {
	return s.client.GetMyInfo()
}

func (s *TelegramAPIService) GetProfileByToken(accessToken string) (*telegramapi.BotInfo, error) {
	// В Telegram Bot API нет токенов доступа пользователей
	// Возвращаем информацию о боте
	return s.client.GetMyInfo()
}

func (s *TelegramAPIService) SendMessage(chatID int64, text string) error {
	req := &telegramapi.SendMessageRequest{
		Text: text,
	}
	_, err := s.client.SendMessage(chatID, req)
	return err
}

func (s *TelegramAPIService) SendMessageToUser(userID int64, message *telegramapi.SendMessageRequest) (*telegramapi.SendMessageResponse, error) {
	return s.client.SendMessageToUser(userID, message)
}

func (s *TelegramAPIService) GetChat(chatID int64) (*telegramapi.Chat, error) {
	return s.client.GetChat(chatID)
}

func (s *TelegramAPIService) GetChatByLink(chatLink string) (*telegramapi.Chat, error) {
	return s.client.GetChatByLink(chatLink)
}

func (s *TelegramAPIService) GetUserInfo(userID int64) (*telegramapi.TelegramUser, error) {
	// В Telegram Bot API нет прямого метода получения информации о пользователе
	// Информация о пользователе приходит в обновлениях
	return nil, fmt.Errorf("not implemented yet")
}

func (s *TelegramAPIService) GetMessages(chatID int64, from, to, count *int64, messageIDs []string) ([]telegramapi.Message, error) {
	return s.client.GetMessages(chatID, from, to, count, messageIDs)
}

func (s *TelegramAPIService) AddMembers(chatID int64, userIDs []int64) error {
	return s.client.AddMembers(chatID, userIDs)
}

func (s *TelegramAPIService) GetChatMembers(chatID int64, marker *int64, count *int, userIDs []int64) (*telegramapi.ChatMembersResponse, error) {
	return s.client.GetChatMembers(chatID, marker, count, userIDs)
}

func (s *TelegramAPIService) EditChat(chatID int64, title *string, icon interface{}) (*telegramapi.Chat, error) {
	return s.client.EditChat(chatID, title, icon)
}

func (s *TelegramAPIService) DeleteChat(chatID int64) error {
	return s.client.DeleteChat(chatID)
}

func (s *TelegramAPIService) RemoveMember(chatID int64, userID int64) error {
	return s.client.RemoveMember(chatID, userID)
}
