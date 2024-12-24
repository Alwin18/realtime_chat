package websocket

import (
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Message struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
}

type clientInfo struct {
	conn     *websocket.Conn
	senderID string
}

type Hub struct {
	clients               sync.Map             // Map senderId to websocket.Conn
	clientRegisterChannel chan *clientInfo     // Register clients
	clientRemoveChannel   chan *websocket.Conn // Remove clients
	broadcastMessage      chan Message         // Broadcast messages
}

func NewHub() *Hub {
	return &Hub{
		clientRegisterChannel: make(chan *clientInfo),
		clientRemoveChannel:   make(chan *websocket.Conn),
		broadcastMessage:      make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case info := <-h.clientRegisterChannel:
			h.clients.Store(info.senderID, info.conn)

		case conn := <-h.clientRemoveChannel:
			h.clients.Range(func(key, value interface{}) bool {
				if value.(*websocket.Conn) == conn {
					h.clients.Delete(key)
					return false
				}
				return true
			})

		case message := <-h.broadcastMessage:
			if value, ok := h.clients.Load(message.ReceiverID); ok {
				conn := value.(*websocket.Conn)
				err := conn.WriteJSON(message)
				if err != nil {
					log.Errorf("Error sending message to %s: %v", message.ReceiverID, err)
					conn.Close()
					h.clients.Delete(message.ReceiverID)
				}
			}
		}
	}
}

func AllowUpgrade(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}

func DirectMessage(hub *Hub) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		defer func() {
			hub.clientRemoveChannel <- conn
			conn.Close()
		}()

		// Get senderId and receiverId from query parameters
		senderID := conn.Query("senderId")
		receiverID := conn.Query("receiverId")
		if senderID == "" || receiverID == "" {
			log.Error("Missing senderId or receiverId")
			conn.Close()
			return
		}

		// Register client
		hub.clientRegisterChannel <- &clientInfo{conn: conn, senderID: senderID}

		// Read messages
		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("Error reading message from %s: %v", senderID, err)
				hub.clientRemoveChannel <- conn
				return
			}

			if messageType == websocket.TextMessage {
				// Broadcast message to the receiver
				message := Message{
					SenderID:   senderID,
					ReceiverID: receiverID,
					Content:    string(msg),
					Timestamp:  time.Now().Format(time.RFC3339),
				}
				hub.broadcastMessage <- message
			}
		}
	}
}
