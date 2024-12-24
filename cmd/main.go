package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/websoket-chat/internal/api"
	websocketHandler "github.com/websoket-chat/internal/websocket"
)

func main() {
	app := fiber.New()

	// Middleware to allow WebSocket upgrades
	app.Use(websocketHandler.AllowUpgrade)

	// Register WebSocket route
	hub := websocketHandler.NewHub()
	go hub.Run()
	app.Get("/ws/direct", websocket.New(websocketHandler.DirectMessage(hub)))

	// Example API route (if needed)
	app.Get("/api/hello", api.HelloHandler)

	log.Info("Server started on port 3000")
	log.Fatal(app.Listen(":3000"))
}
