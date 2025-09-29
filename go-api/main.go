package main

import (
	"log"

	"github.com/DevaSinha/StreamSight/go-api/config"
	"github.com/DevaSinha/StreamSight/go-api/routes"
	"github.com/DevaSinha/StreamSight/go-api/websocket"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDatabase()
	defer config.CloseDatabase()

	hub := websocket.NewHub()
	go hub.Run()

	router := gin.Default()
	routes.SetupRoutes(router, hub)

	log.Println("Server starting on :8080")
	router.Run(":8080")
}
