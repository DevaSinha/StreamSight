package handlers

import (
	"net/http"

	wsHub "github.com/DevaSinha/StreamSight/go-api/websocket"
	"github.com/gin-gonic/gin"
	gorillaWs "github.com/gorilla/websocket"
)

var upgrader = gorillaWs.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeAlertsWS(hub *wsHub.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		hub.Register(conn)

		defer hub.Unregister(conn)

		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}
}
