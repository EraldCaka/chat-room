package main

import (
	"github.com/EraldCaka/chat-room/db"
	"github.com/EraldCaka/chat-room/internal/handlers"
	"github.com/EraldCaka/chat-room/internal/repository"
	"github.com/EraldCaka/chat-room/internal/services"
	"github.com/EraldCaka/chat-room/internal/ws"
	"github.com/EraldCaka/chat-room/router"
	"github.com/EraldCaka/chat-room/util"
	"log"
)

func main() {
	util.InitEnvironmentVariables()
	dbConn, err := db.NewDatabase()
	log.Println(dbConn.GetDB())
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRepository := repository.NewRepository(dbConn.GetDB())
	userService := services.NewService(userRepository)
	userHandler := handlers.NewHandler(userService)
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()
	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:5555")
}
