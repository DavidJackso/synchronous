package v1

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rnegic/synchronous/internal/interfaces"
	"github.com/rnegic/synchronous/pkg/maxapi"
)

type WebhookHandler struct {
	*BaseHandler
	sessionService interfaces.SessionService
	maxAPIService  interfaces.MaxAPIService
}

const welcomeMessage = "Привет! Я Synchronous Bot — помогаю проводить фокус-сессии и синхронно работать с командой.\n\nВот что я умею:\n- запускать solo и групповые сессии по Помодоро с гибкими циклами\n- собирать задачи и отслеживать их выполнение в реальном времени\n- приглашать коллег по ссылке или через Max-чаты\n- сохранять отчёты по каждой сессии и делиться ими\n- создавать обсуждения прямо в Max после завершения фокуса\n\nЧтобы стартовать, просто открой WebApp Synchronous и создай первую сессию — я подскажу каждый шаг."

func NewWebhookHandler(baseHandler *BaseHandler, sessionService interfaces.SessionService, maxAPIService interfaces.MaxAPIService) *WebhookHandler {
	return &WebhookHandler{
		BaseHandler:    baseHandler,
		sessionService: sessionService,
		maxAPIService:  maxAPIService,
	}
}

func (h *WebhookHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Webhook endpoint (публичный, без аутентификации)
	router.POST("/webhook/max", h.handleWebhook)
}

// handleWebhook обрабатывает webhook от Max API
func (h *WebhookHandler) handleWebhook(c *gin.Context) {
	// Читаем тело запроса
	body, err := c.GetRawData()
	if err != nil {
		log.Printf("[Webhook] Failed to read request body: %v", err)
		h.ErrorResponse(c, http.StatusBadRequest, "failed to read request body")
		return
	}

	// Парсим обновление
	update, err := maxapi.ParseUpdate(body)
	if err != nil {
		log.Printf("[Webhook] Failed to parse update: %v", err)
		h.ErrorResponse(c, http.StatusBadRequest, "failed to parse update")
		return
	}

	// Обрабатываем обновление в зависимости от типа
	switch u := update.(type) {
	case *maxapi.MessageCreatedUpdate:
		log.Printf("[Webhook] Received message from user=%d chat=%d text=%q",
			u.Message.Sender.UserID, u.Message.Recipient.ChatID, u.Message.Body.Text)

		if err := h.handleMessageCreated(u); err != nil {
			log.Printf("[Webhook] Failed to handle message_created: %v", err)
			h.ErrorResponse(c, http.StatusInternalServerError, "failed to process message")
			return
		}

		h.SuccessResponse(c, http.StatusOK, gin.H{"status": "processed"})

	case *maxapi.MessageChatCreatedUpdate:
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

func (h *WebhookHandler) handleMessageCreated(update *maxapi.MessageCreatedUpdate) error {
	if update == nil || h.maxAPIService == nil {
		return nil
	}

	text := strings.TrimSpace(update.Message.Body.Text)
	if text == "" {
		return nil
	}

	lowered := strings.ToLower(text)
	if lowered != "/start" && lowered != "start" && lowered != "привет" {
		return nil
	}

	if update.Message.Sender.UserID == 0 {
		return nil
	}

	_, err := h.maxAPIService.SendMessageToUser(update.Message.Sender.UserID, &maxapi.SendMessageRequest{
		Text: welcomeMessage,
	})
	return err
}
