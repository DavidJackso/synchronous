package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
	"github.com/rnegic/synchronous/pkg/jwt"
)

type AuthHandler struct {
	*BaseHandler
	authService  interfaces.AuthService
	tokenManager *jwt.TokenManager
}

func NewAuthHandler(baseHandler *BaseHandler, authService interfaces.AuthService, tokenManager *jwt.TokenManager) *AuthHandler {
	return &AuthHandler{
		BaseHandler:  baseHandler,
		authService:  authService,
		tokenManager: tokenManager,
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
	var req entity.TelegramAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	fmt.Printf("[Auth Handler] 📥 Login request received\n")
	fmt.Printf("[Auth Handler]   Request origin: %s\n", c.Request.Header.Get("Origin"))
	fmt.Printf("[Auth Handler]   User-Agent: %s\n", c.Request.Header.Get("User-Agent"))
	fmt.Printf("[Auth Handler]   X-Forwarded-Proto: %s\n", c.Request.Header.Get("X-Forwarded-Proto"))

	tokens, user, err := h.authService.Login(req.InitData, req.DeviceID)
	if err != nil {
		fmt.Printf("[Auth Handler] ❌ Login failed: %v\n", err)
		h.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	fmt.Printf("[Auth Handler] ✅ Login successful for user: %s (ID: %s)\n", user.Name, user.ID)

	// Set HTTP-only cookies instead of returning tokens in body
	accessTTL := h.tokenManager.GetAccessTTL()
	refreshTTL := h.tokenManager.GetRefreshTTL()

	fmt.Printf("[Auth Handler] 🍪 Setting cookies (access TTL: %ds, refresh TTL: %ds)\n", accessTTL, refreshTTL)
	h.setAccessTokenCookie(c, tokens.AccessToken, accessTTL)
	h.setRefreshTokenCookie(c, tokens.RefreshToken, refreshTTL)

	// Return only user data (no tokens in response body)
	h.SuccessResponse(c, http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"avatarUrl": user.AvatarURL,
			"createdAt": user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})

	fmt.Printf("[Auth Handler] ✅ Response sent with cookies\n")
}

func (h *AuthHandler) refresh(c *gin.Context) {
	// Read refresh token from cookie instead of request body
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		h.ErrorResponse(c, http.StatusUnauthorized, "refresh token required")
		return
	}

	tokens, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		h.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set new access token cookie
	accessTTL := h.tokenManager.GetAccessTTL()
	h.setAccessTokenCookie(c, tokens.AccessToken, accessTTL)

	// Optionally update refresh token cookie if a new one was generated
	if tokens.RefreshToken != "" {
		refreshTTL := h.tokenManager.GetRefreshTTL()
		h.setRefreshTokenCookie(c, tokens.RefreshToken, refreshTTL)
	}

	// Return success response (no tokens in body)
	h.SuccessResponse(c, http.StatusOK, gin.H{
		"success": true,
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

	// Clear cookies by setting them with MaxAge=0
	h.clearAccessTokenCookie(c)
	h.clearRefreshTokenCookie(c)

	c.Status(http.StatusNoContent)
}
