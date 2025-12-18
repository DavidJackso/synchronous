package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rnegic/synchronous/internal/interfaces"
)

func AuthMiddleware(authService interfaces.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read access token from cookie instead of Authorization header
		accessToken, err := c.Cookie("access_token")
		if err != nil || accessToken == "" {
			fmt.Printf("[Auth Middleware] ❌ No access_token cookie found: %v\n", err)
			fmt.Printf("[Auth Middleware]   Request URL: %s\n", c.Request.URL.Path)
			fmt.Printf("[Auth Middleware]   Request Method: %s\n", c.Request.Method)
			fmt.Printf("[Auth Middleware]   All cookies: %v\n", c.Request.Cookies())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		fmt.Printf("[Auth Middleware] ✅ Found access_token cookie (length: %d)\n", len(accessToken))

		userID, err := authService.ValidateToken(accessToken)
		if err != nil {
			fmt.Printf("[Auth Middleware] ❌ Token validation failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		fmt.Printf("[Auth Middleware] ✅ Token validated for user: %s\n", userID)

		// Сохраняем userID в контексте gin
		c.Set("userID", userID)
		c.Set("X-User-ID", userID) // Для обратной совместимости

		c.Next()
	}
}
