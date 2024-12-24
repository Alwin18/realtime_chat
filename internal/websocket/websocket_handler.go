package websocket

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/websoket-chat/internal/model"
	"github.com/websoket-chat/internal/repository"
	"github.com/websoket-chat/utils"

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
	chatMessageRepository *repository.ChatMessageRepository
}

func NewHub(chatMessageRepository *repository.ChatMessageRepository) *Hub {
	return &Hub{
		clientRegisterChannel: make(chan *clientInfo),
		clientRemoveChannel:   make(chan *websocket.Conn),
		broadcastMessage:      make(chan Message),
		chatMessageRepository: chatMessageRepository,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case info := <-h.clientRegisterChannel:
			// Remove the previous connection with the same senderId if it exists
			if oldConn, ok := h.clients.Load(info.conn.Query("senderId", "")); ok {
				log.Infof("Closing old connection for senderId %s and receiverId %s", info.conn.Query("senderId", ""), oldConn.(*websocket.Conn).Query("receiverId", ""))
				oldConn.(*websocket.Conn).Close()
				h.clients.Delete(info.conn.Query("senderId", ""))
			}

			// Store the new connection by senderId
			h.clients.Store(info.conn.Query("senderId", ""), info.conn)
			log.Infof("New connection for senderId %s and receiverId %s", info.conn.Query("senderId", ""), info.conn.Query("receiverId", ""))

		case conn := <-h.clientRemoveChannel:
			h.clients.Range(func(key, value interface{}) bool {
				if value.(*websocket.Conn) == conn {
					log.Info("Closing connection for senderId ", key)
					h.clients.Delete(key)
					return false
				}
				return true
			})

		case message := <-h.broadcastMessage:
			// Broadcast message to specific receiver
			if value, ok := h.clients.Load(message.ReceiverID); ok {
				err := value.(*websocket.Conn).WriteJSON(message)
				if err != nil {
					log.Errorf("Error sending message to %s: %v", message.ReceiverID, err)
					value.(*websocket.Conn).Close()
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
		if conn.Query("senderId") == "" || conn.Query("receiverId", "") == "" {
			log.Error("Missing senderId or receiverId")
			conn.Close()
			return
		}

		// Register client
		hub.clientRegisterChannel <- &clientInfo{
			conn:     conn,
			senderID: conn.Query("senderId"),
		}

		// Read messages
		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("Error reading message from %s: %v", conn.Query("senderId"), err)
				hub.clientRemoveChannel <- conn
				return
			}

			if messageType == websocket.TextMessage {
				// Save message to database
				senderId, err := strconv.ParseInt(conn.Query("senderId"), 10, 64)
				if err != nil {
					log.Errorf("Error parsing senderId: %v", err)
				}

				receiverId, err := strconv.ParseInt(conn.Query("receiverId"), 10, 64)
				if err != nil {
					log.Errorf("Error parsing senderId: %v", err)
				}

				err = hub.chatMessageRepository.SaveMessage(model.ChatMessage{
					SenderID:   senderId,
					ReceiverID: receiverId,
					Content:    string(msg),
					Timestamp:  utils.ConvertToJakartaTime(time.Now()),
				})
				if err != nil {
					log.Errorf("Error saving message to DB: %v", err)
				}

				// broadcast the message to specific receiver
				hub.broadcastMessage <- Message{
					SenderID:   conn.Query("senderId", ""),
					ReceiverID: conn.Query("receiverId", ""),
					Content:    string(msg),
					Timestamp:  utils.ConvertToJakartaTime(time.Now()).Format(time.RFC3339),
				}
			}
		}
	}
}
