package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rnegic/synchronous/internal/interfaces"
)

func AuthMiddleware(authService interfaces.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read access token from cookie instead of Authorization header
		accessToken, err := c.Cookie("access_token")
		if err != nil || accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		userID, err := authService.ValidateToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		// Сохраняем userID в контексте gin
		c.Set("userID", userID)
		c.Set("X-User-ID", userID) // Для обратной совместимости

		c.Next()
	}
}
