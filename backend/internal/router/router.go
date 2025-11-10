package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	// Устанавливаем режим gin (release для продакшена)
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// CORS configuration for HTTP-only cookies
	config := cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000", // Local development
			"http://focus-sync.ru",  // Production HTTP
			"https://focus-sync.ru", // Production HTTPS
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept",
		},
		AllowCredentials: true, // Critical for cookies
		MaxAge:           12 * time.Hour,
	}

	// Добавляем middleware
	router.Use(cors.New(config))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}
