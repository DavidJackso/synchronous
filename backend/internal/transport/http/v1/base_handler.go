package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// Helper методы для работы с gin.Context
func (h *BaseHandler) GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func (h *BaseHandler) ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

func (h *BaseHandler) SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// isSecureRequest checks if the request is over HTTPS
// Checks X-Forwarded-Proto header (set by nginx) or request scheme
func (h *BaseHandler) isSecureRequest(c *gin.Context) bool {
	// Check X-Forwarded-Proto header (set by reverse proxy)
	if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" {
		return true
	}
	// Check request scheme directly
	return c.Request.TLS != nil || c.Request.URL.Scheme == "https"
}

// setAccessTokenCookie sets the access token as an HTTP-only cookie
func (h *BaseHandler) setAccessTokenCookie(c *gin.Context, token string, maxAge int) {
	secure := h.isSecureRequest(c)
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)
}

// setRefreshTokenCookie sets the refresh token as an HTTP-only cookie
func (h *BaseHandler) setRefreshTokenCookie(c *gin.Context, token string, maxAge int) {
	secure := h.isSecureRequest(c)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/api/v1/auth/refresh",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)
}

// clearAccessTokenCookie clears the access token cookie
func (h *BaseHandler) clearAccessTokenCookie(c *gin.Context) {
	secure := h.isSecureRequest(c)
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)
}

// clearRefreshTokenCookie clears the refresh token cookie
func (h *BaseHandler) clearRefreshTokenCookie(c *gin.Context) {
	secure := h.isSecureRequest(c)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(c.Writer, cookie)
}
