package maxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

type BotInfo struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Name      string `json:"name,omitempty"`
	Username  string `json:"username,omitempty"`
	IsBot     bool   `json:"is_bot"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type MaxUser struct {
	UserID    int64  `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	IsBot     bool   `json:"is_bot"`
	AvatarURL string `json:"avatar_url,omitempty"`
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

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func NewClient(baseURL, accessToken string) *Client {
	return &Client{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetMyInfo() (*BotInfo, error) {
	req, err := c.newRequest("GET", "/me", nil)
	if err != nil {
		return nil, err
	}

	var botInfo BotInfo
	if err := c.doRequest(req, &botInfo); err != nil {
		return nil, err
	}

	return &botInfo, nil
}

func (c *Client) SendMessage(chatID int64, message *SendMessageRequest) (*SendMessageResponse, error) {
	endpoint := fmt.Sprintf("/messages?chat_id=%d", chatID)
	req, err := c.newRequest("POST", endpoint, message)
	if err != nil {
		return nil, err
	}

	var response SendMessageResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SendMessageToUser отправляет сообщение пользователю (в личный чат с ботом)
func (c *Client) SendMessageToUser(userID int64, message *SendMessageRequest) (*SendMessageResponse, error) {
	endpoint := fmt.Sprintf("/messages?user_id=%d", userID)
	req, err := c.newRequest("POST", endpoint, message)
	if err != nil {
		return nil, err
	}

	var response SendMessageResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetChat(chatID int64) (*Chat, error) {
	endpoint := fmt.Sprintf("/chats/%d", chatID)
	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var chat Chat
	if err := c.doRequest(req, &chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (c *Client) GetChatByLink(chatLink string) (*Chat, error) {
	endpoint := fmt.Sprintf("/chats/%s", url.PathEscape(chatLink))
	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var chat Chat
	if err := c.doRequest(req, &chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

// GetMessages получает сообщения из чата
func (c *Client) GetMessages(chatID int64, from, to, count *int64, messageIDs []string) ([]Message, error) {
	endpoint := fmt.Sprintf("/messages?chat_id=%d", chatID)

	q := url.Values{}
	if from != nil {
		q.Set("from", fmt.Sprintf("%d", *from))
	}
	if to != nil {
		q.Set("to", fmt.Sprintf("%d", *to))
	}
	if count != nil {
		q.Set("count", fmt.Sprintf("%d", *count))
	}
	if len(messageIDs) > 0 {
		q.Set("message_ids", fmt.Sprintf("%v", messageIDs))
	}

	if len(q) > 0 {
		endpoint += "&" + q.Encode()
	}

	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	type MessageListResponse struct {
		Messages []Message `json:"messages"`
	}

	var response MessageListResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Messages, nil
}

// AddMembers добавляет участников в чат
func (c *Client) AddMembers(chatID int64, userIDs []int64) error {
	endpoint := fmt.Sprintf("/chats/%d/members", chatID)

	reqBody := struct {
		UserIDs []int64 `json:"user_ids"`
	}{
		UserIDs: userIDs,
	}

	req, err := c.newRequest("POST", endpoint, reqBody)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// GetChatMembers получает список участников чата
func (c *Client) GetChatMembers(chatID int64, marker *int64, count *int, userIDs []int64) (*ChatMembersResponse, error) {
	endpoint := fmt.Sprintf("/chats/%d/members", chatID)

	q := url.Values{}
	if marker != nil {
		q.Set("marker", fmt.Sprintf("%d", *marker))
	}
	if count != nil {
		q.Set("count", fmt.Sprintf("%d", *count))
	}
	if len(userIDs) > 0 {
		// Max API принимает user_ids как массив через query параметр
		for _, id := range userIDs {
			q.Add("user_ids", fmt.Sprintf("%d", id))
		}
	}

	if len(q) > 0 {
		endpoint += "&" + q.Encode()
	}

	req, err := c.newRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var response ChatMembersResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
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

// EditChat редактирует информацию о чате
func (c *Client) EditChat(chatID int64, title *string, icon interface{}) (*Chat, error) {
	endpoint := fmt.Sprintf("/chats/%d", chatID)

	reqBody := make(map[string]interface{})
	if title != nil {
		reqBody["title"] = *title
	}
	if icon != nil {
		reqBody["icon"] = icon
	}

	req, err := c.newRequest("PATCH", endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	var chat Chat
	if err := c.doRequest(req, &chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

// DeleteChat удаляет чат (только если бот является владельцем или администратором)
func (c *Client) DeleteChat(chatID int64) error {
	endpoint := fmt.Sprintf("/chats/%d", chatID)
	req, err := c.newRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

// RemoveMember удаляет участника из чата
func (c *Client) RemoveMember(chatID int64, userID int64) error {
	endpoint := fmt.Sprintf("/chats/%d/members/%d", chatID, userID)
	req, err := c.newRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	return c.doRequest(req, nil)
}

func (c *Client) newRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Добавляем access_token в query параметры
	q := req.URL.Query()
	q.Set("access_token", c.accessToken)
	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (c *Client) doRequest(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return fmt.Errorf("API error [%d]: %s - %s", resp.StatusCode, errResp.Code, errResp.Message)
		}
		return fmt.Errorf("API error [%d]: %s", resp.StatusCode, resp.Status)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
