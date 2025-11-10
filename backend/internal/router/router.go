package router

import (
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	// Устанавливаем режим gin (release для продакшена)
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Добавляем middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}
