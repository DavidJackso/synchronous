package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type UserHandler struct {
	*BaseHandler
	userService interfaces.UserService
}

func NewUserHandler(baseHandler *BaseHandler, userService interfaces.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: baseHandler,
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/me", h.getMe)
		users.GET("/contacts", h.getContacts)
	}
}

func (h *UserHandler) getMe(c *gin.Context) {
	userID := h.GetUserID(c)
	if userID == "" {
		h.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, stats, err := h.userService.GetProfile(userID)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.SuccessResponse(c, http.StatusOK, gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"avatarUrl": user.AvatarURL,
		"stats": gin.H{
			"totalSessions":  stats.TotalSessions,
			"totalFocusTime": stats.TotalFocusTime,
			"currentStreak":  stats.CurrentStreak,
		},
		"createdAt": user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *UserHandler) getContacts(c *gin.Context) {
	userID := h.GetUserID(c)
	if userID == "" {
		h.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	contacts, err := h.userService.GetContacts(userID)
	if err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	contactsList := make([]gin.H, 0, len(contacts))
	for _, contact := range contacts {
		contactsList = append(contactsList, gin.H{
			"id":           contact.ID,
			"name":         contact.Name,
			"avatarUrl":    contact.AvatarURL,
			"isRegistered": true, // В реальности нужно проверить
		})
	}

	h.SuccessResponse(c, http.StatusOK, gin.H{
		"contacts": contactsList,
	})
}
