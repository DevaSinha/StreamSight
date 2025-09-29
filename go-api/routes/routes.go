package routes

import (
	"github.com/DevaSinha/StreamSight/go-api/handlers"
	"github.com/DevaSinha/StreamSight/go-api/middleware"
	wsHub "github.com/DevaSinha/StreamSight/go-api/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, hub *wsHub.Hub) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/ws/alerts", handlers.ServeAlertsWS(hub))

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/me", middleware.AuthMiddleware(), handlers.GetMe)
	}

	// Camera routes
	cameras := router.Group("/cameras")
	cameras.Use(middleware.AuthMiddleware())
	{
		cameras.GET("", handlers.ListCameras)
		cameras.POST("", handlers.CreateCamera)
		cameras.GET("/:id", handlers.GetCamera)
		cameras.PUT("/:id", handlers.UpdateCamera)
		cameras.DELETE("/:id", handlers.DeleteCamera)
	}
}
