package telegramapi

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	bot *tgbotapi.BotAPI
}

type BotInfo struct {
	UserID        int64  `json:"user_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name,omitempty"`
	Name          string `json:"name,omitempty"`
	Username      string `json:"username,omitempty"`
	IsBot         bool   `json:"is_bot"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	FullAvatarURL string `json:"full_avatar_url,omitempty"`
	Description   string `json:"description,omitempty"`
}

type TelegramUser struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	IsBot     bool   `json:"is_bot"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Name      string `json:"name,omitempty"`
}

type Chat struct {
	ChatID            int64  `json:"chat_id"`
	Type              string `json:"type"`
	Status            string `json:"status"`
	Title             string `json:"title,omitempty"`
	LastEventTime     int64  `json:"last_event_time"`
	ParticipantsCount int    `json:"participants_count"`
}

type Message struct {
	Sender struct {
		UserID    int64  `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name,omitempty"`
		Username  string `json:"username,omitempty"`
	} `json:"sender"`
	Recipient struct {
		ChatID   int64  `json:"chat_id"`
		UserID   int64  `json:"user_id"`
		ChatType string `json:"chat_type"`
	} `json:"recipient"`
	Timestamp int64 `json:"timestamp"`
	Body      struct {
		Mid         string        `json:"mid"`
		Text        string        `json:"text,omitempty"`
		Attachments []interface{} `json:"attachments,omitempty"`
	} `json:"body"`
}

type SendMessageRequest struct {
	Text        string        `json:"text,omitempty"`
	Attachments []interface{} `json:"attachments,omitempty"`
}

type SendMessageResponse struct {
	Message Message `json:"message"`
}

type ChatMember struct {
	UserID         int64    `json:"user_id"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name,omitempty"`
	Username       string   `json:"username,omitempty"`
	IsBot          bool     `json:"is_bot"`
	AvatarURL      string   `json:"avatar_url,omitempty"`
	LastAccessTime int64    `json:"last_access_time"`
	IsOwner        bool     `json:"is_owner"`
	IsAdmin        bool     `json:"is_admin"`
	JoinTime       int64    `json:"join_time"`
	Permissions    []string `json:"permissions,omitempty"`
	Alias          string   `json:"alias,omitempty"`
}

type ChatMembersResponse struct {
	Members []ChatMember `json:"members"`
	Marker  *int64       `json:"marker,omitempty"`
}

func NewClient(botToken string) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &Client{
		bot: bot,
	}, nil
}

func (c *Client) GetMyInfo() (*BotInfo, error) {
	me, err := c.bot.GetMe()
	if err != nil {
		return nil, err
	}

	return &BotInfo{
		UserID:    int64(me.ID),
		FirstName: me.FirstName,
		LastName:  me.LastName,
		Username:  me.UserName,
		IsBot:     me.IsBot,
		Name:      me.FirstName + " " + me.LastName,
	}, nil
}

func (c *Client) SendMessage(chatID int64, message *SendMessageRequest) (*SendMessageResponse, error) {
	msg := tgbotapi.NewMessage(chatID, message.Text)
	if message != nil && message.Text != "" {
		msg.Text = message.Text
	}

	sentMsg, err := c.bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &SendMessageResponse{
		Message: convertMessage(sentMsg),
	}, nil
}

func (c *Client) SendMessageToUser(userID int64, message *SendMessageRequest) (*SendMessageResponse, error) {
	// В Telegram отправка сообщения пользователю по ID - это отправка в личный чат
	msg := tgbotapi.NewMessage(userID, message.Text)
	if message != nil && message.Text != "" {
		msg.Text = message.Text
	}

	sentMsg, err := c.bot.Send(msg)
	if err != nil {
		return nil, err
	}

	return &SendMessageResponse{
		Message: convertMessage(sentMsg),
	}, nil
}

func (c *Client) GetChat(chatID int64) (*Chat, error) {
	chatConfig := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	}
	chat, err := c.bot.GetChat(chatConfig)
	if err != nil {
		return nil, err
	}

	return convertChat(&chat), nil
}

func (c *Client) GetChatByLink(chatLink string) (*Chat, error) {
	// В Telegram нужно получить чат по invite link
	// Для этого нужно сначала принять приглашение или использовать другие методы
	// Пока возвращаем ошибку, так как это требует дополнительной логики
	return nil, nil
}

func (c *Client) GetMessages(chatID int64, from, to, count *int64, messageIDs []string) ([]Message, error) {
	// Telegram Bot API не поддерживает получение истории сообщений напрямую
	// Сообщения приходят через webhook/updates
	// Возвращаем пустой список
	return []Message{}, nil
}

func (c *Client) AddMembers(chatID int64, userIDs []int64) error {
	// В Telegram боты не могут добавлять участников в группы напрямую
	// Нужны специальные права администратора
	// Пока возвращаем nil
	return nil
}

func (c *Client) GetChatMembers(chatID int64, marker *int64, count *int, userIDs []int64) (*ChatMembersResponse, error) {
	members, err := c.bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	})
	if err != nil {
		return nil, err
	}

	result := &ChatMembersResponse{
		Members: make([]ChatMember, 0, len(members)),
	}

	for _, member := range members {
		result.Members = append(result.Members, convertChatMember(member))
	}

	// Telegram Bot API не возвращает общее количество участников напрямую
	// Для этого нужно использовать другие методы API

	return result, nil
}

func (c *Client) EditChat(chatID int64, title *string, icon interface{}) (*Chat, error) {
	if title != nil {
		config := tgbotapi.NewChatTitle(chatID, *title)
		_, err := c.bot.Request(config)
		if err != nil {
			return nil, err
		}
	}

	chatConfig := tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: chatID,
		},
	}
	chat, err := c.bot.GetChat(chatConfig)
	if err != nil {
		return nil, err
	}

	return convertChat(&chat), nil
}

func (c *Client) DeleteChat(chatID int64) error {
	// В Telegram боты не могут удалять чаты
	// Могут только покинуть группу
	config := tgbotapi.LeaveChatConfig{
		ChatID: chatID,
	}
	_, err := c.bot.Request(config)
	return err
}

func (c *Client) RemoveMember(chatID int64, userID int64) error {
	config := tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
	}
	_, err := c.bot.Request(config)
	return err
}

// --- Helpers ----------------------------------------------------------------

func convertChat(chat *tgbotapi.Chat) *Chat {
	if chat == nil {
		return nil
	}

	chatType := string(chat.Type)
	status := "active"
	if chat.Permissions != nil {
		status = "active"
	}

	// В Telegram Bot API количество участников не всегда доступно
	// Можно получить через GetChatMembersCount, но это требует отдельного запроса
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

func convertChatMember(member tgbotapi.ChatMember) ChatMember {
	user := member.User
	isAdmin := member.Status == "administrator" || member.Status == "creator"
	isOwner := member.Status == "creator"

	return ChatMember{
		UserID:         int64(user.ID),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Username:       user.UserName,
		IsBot:          user.IsBot,
		LastAccessTime: 0,
		IsOwner:        isOwner,
		IsAdmin:        isAdmin,
		JoinTime:       0,
	}
}

func convertMessage(msg tgbotapi.Message) Message {
	var result Message

	result.Sender.UserID = int64(msg.From.ID)
	result.Sender.FirstName = msg.From.FirstName
	result.Sender.LastName = msg.From.LastName
	result.Sender.Username = msg.From.UserName

	result.Recipient.ChatID = msg.Chat.ID
	result.Recipient.UserID = int64(msg.From.ID)
	result.Recipient.ChatType = string(msg.Chat.Type)

	result.Timestamp = int64(msg.Date)
	result.Body.Mid = ""
	if msg.MessageID != 0 {
		result.Body.Mid = string(rune(msg.MessageID))
	}
	result.Body.Text = msg.Text

	return result
}
