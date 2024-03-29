package router

import (
	"github.com/EraldCaka/chat-room/internal/handlers"
	"github.com/EraldCaka/chat-room/internal/middleware"
	"github.com/EraldCaka/chat-room/internal/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

var r *gin.Engine

func InitRouter(userHandler *handlers.Handler, wsHandler *ws.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5174"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.Use(middleware.AuthMiddleware())

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)
	r.GET("/cookie", userHandler.GetJWTFromCookie)

	r.POST("/ws/createRoom", wsHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	r.GET("/ws/getRooms", wsHandler.GetRooms)
	r.GET("/ws/room/clients/:roomID", wsHandler.GetRoomActiveClients)
	r.GET("/ws/client/close/:userID/:roomID", wsHandler.CloseWSConnection)

}

func Start(addr string) error {
	return r.Run(addr)
}
