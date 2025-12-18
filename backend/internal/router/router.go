package router

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	// Устанавливаем режим gin (release для продакшена)
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// CORS configuration for HTTP-only cookies
	allowedOrigins := []string{
		"http://localhost:3000",    // Local development
		"http://focus-sync.ru",     // Production HTTP
		"https://focus-sync.ru",    // Production HTTPS
		"http://tg.focus-sync.ru",  // Subdomain HTTP
		"https://tg.focus-sync.ru", // Subdomain HTTPS
	}

	config := cors.Config{
		// Don't use AllowOrigins when using AllowOriginFunc
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length", "Set-Cookie",
		},
		AllowCredentials: true, // Critical for cookies
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			// Log all origins for debugging
			fmt.Printf("[CORS] 🔍 Checking origin: %s\n", origin)
			// Check against allowed origins
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					fmt.Printf("[CORS] ✅ Origin allowed: %s\n", origin)
					return true
				}
			}
			fmt.Printf("[CORS] ❌ Origin not allowed: %s\n", origin)
			return false
		},
	}

	// Добавляем middleware для логирования cookies (до CORS)
	router.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		cookieHeader := c.GetHeader("Cookie")
		fmt.Printf("[Request] 📥 %s %s\n", c.Request.Method, c.Request.URL.Path)
		fmt.Printf("[Request]   Origin: %s\n", origin)
		if cookieHeader != "" {
			fmt.Printf("[Request] 🍪 Cookies received: %s\n", cookieHeader)
		} else {
			fmt.Printf("[Request] ❌ No Cookie header\n")
		}
		c.Next()
	})

	// Добавляем middleware
	router.Use(cors.New(config))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}
