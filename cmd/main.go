package main

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/websoket-chat/config"
	"github.com/websoket-chat/internal/api"
	"github.com/websoket-chat/internal/repository"
	websocketHandler "github.com/websoket-chat/internal/websocket"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Warnf("Error loading .env file")
	}

	// Load application configuration
	cfg := config.LoadConfig()

	// Connect to the database
	db := config.ConnectDatabase(cfg)
	if db == nil {
		log.Fatal("Failed to connect to database")
		return
	}

	// Migrate database tables
	err = config.MigrateTable(db)
	if err != nil {
		log.Fatal("Error migrating tables: ", err)
		return
	}

	// Initialize Fiber app
	app := fiber.New()

	// Initialize the repository for chat messages
	chatMessageRepo := repository.NewChatMessageRepository(db)

	// Register WebSocket route with AllowUpgrade middleware
	hub := websocketHandler.NewHub(chatMessageRepo)
	go hub.Run()

	// WebSocket route for chat
	app.Get("/ws/direct", websocketHandler.AllowUpgrade, websocket.New(websocketHandler.DirectMessage(hub)))

	// Initialize API handler
	apiHandler := api.NewApiHandler(chatMessageRepo)

	// API route for chat history
	v1 := app.Group("/api/v1")
	v1.Get("hello", api.HelloHandler)
	v1.Get("chat-history", apiHandler.GetChatHistory)

	// Start the server
	log.Info("Server started on port 3000")
	log.Fatal(app.Listen(":3000"))
}
