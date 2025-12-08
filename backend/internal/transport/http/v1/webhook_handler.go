package v1

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rnegic/synchronous/internal/interfaces"
	"github.com/rnegic/synchronous/pkg/telegramapi"
)

type WebhookHandler struct {
	*BaseHandler
	sessionService     interfaces.SessionService
	telegramAPIService interfaces.TelegramAPIService
	authService        interfaces.AuthService
}

const welcomeMessage = `👋 Привет! Это бот Синхрон - твой помощник для фокус-сессий и синхронной работы с командой.

🚀 Что я умею:
• Запускать одиночные и групповые сессии по Помодоро с гибкими циклами
• Собирать задачи и отслеживать их выполнение в реальном времени
• Приглашать коллег по ссылке
• Сохранять отчёты по каждой сессии и делиться ими

📱 Чтобы начать работу:
1. Открой веб-приложение Синхрона
2. Создай свою первую сессию
3. Я подскажу каждый шаг!

💡 Команды:
/start - показать это приветствие
/restart - сбросить авторизацию и начать заново`

const restartMessage = `🔄 Авторизация сброшена!

Теперь тебе нужно:
1. Открой веб-приложение Синхрона
2. Авторизуйся заново через Telegram
3. Начни новую сессию!

Если возникли проблемы, напиши /start для получения помощи.`

func NewWebhookHandler(baseHandler *BaseHandler, sessionService interfaces.SessionService, telegramAPIService interfaces.TelegramAPIService, authService interfaces.AuthService) *WebhookHandler {
	return &WebhookHandler{
		BaseHandler:        baseHandler,
		sessionService:     sessionService,
		telegramAPIService: telegramAPIService,
		authService:        authService,
	}
}

func (h *WebhookHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Webhook endpoint (публичный, без аутентификации)
	router.POST("/webhook/telegram", h.handleWebhook)
}

// handleWebhook обрабатывает webhook от Telegram Bot API
func (h *WebhookHandler) handleWebhook(c *gin.Context) {
	// Читаем тело запроса
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("[Webhook] ❌ Failed to read request body: %v", err)
		h.ErrorResponse(c, http.StatusBadRequest, "failed to read request body")
		return
	}

	log.Printf("[Webhook] 📥 Received webhook, body length: %d bytes", len(body))

	// Парсим обновление
	update, err := telegramapi.ParseUpdate(body)
	if err != nil {
		log.Printf("[Webhook] ❌ Failed to parse update: %v", err)
		log.Printf("[Webhook] Raw body (first 500 chars): %.500s", string(body))
		h.ErrorResponse(c, http.StatusBadRequest, "failed to parse update")
		return
	}

	log.Printf("[Webhook] ✅ Parsed update type: %T", update)

	// Обрабатываем обновление в зависимости от типа
	switch u := update.(type) {
	case *telegramapi.MessageCreatedUpdate:
		log.Printf("[Webhook] 📨 Received message from user=%d chat=%d text=%q",
			u.Message.Sender.UserID, u.Message.Recipient.ChatID, u.Message.Body.Text)

		if err := h.handleMessageCreated(u); err != nil {
			log.Printf("[Webhook] ❌ Failed to handle message_created: %v", err)
			h.ErrorResponse(c, http.StatusInternalServerError, "failed to process message")
			return
		}

		log.Printf("[Webhook] ✅ Message processed successfully")
		h.SuccessResponse(c, http.StatusOK, gin.H{"status": "processed"})

	case *telegramapi.MessageChatCreatedUpdate:
		// Обрабатываем создание чата
		log.Printf("[Webhook] Received chat created update: chatID=%d, startPayload=%s",
			u.Chat.ChatID, u.StartPayload)

		if err := h.sessionService.HandleChatCreated(update); err != nil {
			log.Printf("[Webhook] Failed to handle chat created: %v", err)
			h.ErrorResponse(c, http.StatusInternalServerError, "failed to process chat creation")
			return
		}

		log.Printf("[Webhook] ✅ Chat created successfully: chatID=%d", u.Chat.ChatID)
		h.SuccessResponse(c, http.StatusOK, gin.H{"status": "processed"})

	default:
		// Другие типы обновлений пока не обрабатываем
		log.Printf("[Webhook] Received unhandled update type: %T", update)
		h.SuccessResponse(c, http.StatusOK, gin.H{"status": "ignored"})
	}
}

func (h *WebhookHandler) handleMessageCreated(update *telegramapi.MessageCreatedUpdate) error {
	if update == nil {
		log.Printf("[Webhook] ⚠️ handleMessageCreated: update is nil")
		return nil
	}

	if h.telegramAPIService == nil {
		log.Printf("[Webhook] ⚠️ handleMessageCreated: telegramAPIService is nil")
		return nil
	}

	text := strings.TrimSpace(update.Message.Body.Text)
	if text == "" {
		log.Printf("[Webhook] ⚠️ handleMessageCreated: empty text, ignoring")
		return nil
	}

	lowered := strings.ToLower(text)
	telegramUserID := update.Message.Sender.UserID

	log.Printf("[Webhook] 🔍 Processing message: text=%q, lowered=%q, userID=%d", text, lowered, telegramUserID)

	if telegramUserID == 0 {
		log.Printf("[Webhook] ⚠️ handleMessageCreated: telegramUserID is 0, ignoring")
		return nil
	}

	// Обработка команды /start
	if lowered == "/start" || lowered == "start" || lowered == "привет" {
		log.Printf("[Webhook] 🚀 Handling /start command for user=%d", telegramUserID)
		_, err := h.telegramAPIService.SendMessageToUser(telegramUserID, &telegramapi.SendMessageRequest{
			Text: welcomeMessage,
		})
		if err != nil {
			log.Printf("[Webhook] ❌ Failed to send welcome message: %v", err)
			return err
		}
		log.Printf("[Webhook] ✅ Welcome message sent to user=%d", telegramUserID)
		return nil
	}

	// Обработка команды /restart
	if lowered == "/restart" || lowered == "restart" {
		log.Printf("[Webhook] 🔄 Handling /restart command for user=%d", telegramUserID)
		return h.handleRestart(telegramUserID)
	}

	log.Printf("[Webhook] ℹ️ Unknown command or message, ignoring: %q", text)
	return nil
}

// handleRestart обрабатывает команду /restart - сбрасывает авторизацию пользователя
func (h *WebhookHandler) handleRestart(telegramUserID int64) error {
	log.Printf("[Webhook] 🔄 Processing /restart command for user=%d", telegramUserID)

	// Получаем пользователя по TelegramUserID
	// Используем userRepo через authService
	user, err := h.authService.GetUserByTelegramID(telegramUserID)
	if err != nil {
		log.Printf("[Webhook] ⚠️ User not found for telegramUserID=%d: %v", telegramUserID, err)
		// Пользователь не найден - все равно отправляем сообщение
		_, sendErr := h.telegramAPIService.SendMessageToUser(telegramUserID, &telegramapi.SendMessageRequest{
			Text: restartMessage,
		})
		return sendErr
	}

	// Выполняем logout для пользователя
	if user != nil && user.ID != "" {
		if err := h.authService.Logout(user.ID); err != nil {
			log.Printf("[Webhook] ⚠️ Failed to logout user=%s: %v", user.ID, err)
		} else {
			log.Printf("[Webhook] ✅ Logout successful for user=%s (telegramUserID=%d)", user.ID, telegramUserID)
		}
	}

	// Отправляем сообщение пользователю
	_, err = h.telegramAPIService.SendMessageToUser(telegramUserID, &telegramapi.SendMessageRequest{
		Text: restartMessage,
	})
	return err
}
