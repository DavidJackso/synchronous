package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

type AuthHandler struct {
	*BaseHandler
	authService interfaces.AuthService
}

func NewAuthHandler(baseHandler *BaseHandler, authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: baseHandler,
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.login)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
	}
}

func (h *AuthHandler) login(c *gin.Context) {
	var req entity.MaxAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, user, err := h.authService.Login(req.MaxToken, req.DeviceID)
	if err != nil {
		h.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.SuccessResponse(c, http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"avatarUrl": user.AvatarURL,
			"createdAt": user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
}

func (h *AuthHandler) refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		h.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.SuccessResponse(c, http.StatusOK, gin.H{
		"accessToken": tokens.AccessToken,
	})
}

func (h *AuthHandler) logout(c *gin.Context) {
	userID := h.GetUserID(c)
	if userID == "" {
		h.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.authService.Logout(userID); err != nil {
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
