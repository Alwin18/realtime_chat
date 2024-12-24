package main

import (
	"fmt"
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

type hub struct {
	clients               sync.Map
	clientRegisterChannel chan *clientInfo
	clientRemoveChannel   chan *websocket.Conn
	broadcastMessage      chan Message
}

func (h *hub) run() {
	for {
		select {
		case info := <-h.clientRegisterChannel:
			// Register client by senderId
			h.clients.Store(info.senderID, info.conn)
		case conn := <-h.clientRemoveChannel:
			// Remove client by connection
			h.clients.Range(func(key, value interface{}) bool {
				if value.(*websocket.Conn) == conn {
					h.clients.Delete(key)
					return false
				}
				return true
			})
		case message := <-h.broadcastMessage:
			// Broadcast message only to the receiver
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

func DirectMessage(hub *hub) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		defer func() {
			hub.clientRemoveChannel <- conn
			conn.Close()
		}()

		senderID := conn.Query("senderId")
		receiverID := conn.Query("receiverId")
		if senderID == "" || receiverID == "" {
			conn.Close()
			return
		}

		// Atur batas waktu membaca pesan
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Tangani frame Pong dari klien
		conn.SetPongHandler(func(appData string) error {
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})

		// Daftarkan koneksi
		hub.clientRegisterChannel <- &clientInfo{conn: conn, senderID: senderID}

		// Kirim frame Ping secara berkala
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Error(fmt.Sprintf("Error sending ping: %v", err))
					hub.clientRemoveChannel <- conn
					return
				}
			}
		}()

		// Baca pesan dari klien
		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Error(fmt.Sprintf("Error reading message: %v", err))
				hub.clientRemoveChannel <- conn
				return
			}

			if messageType == websocket.TextMessage {
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

func main() {
	h := &hub{
		clientRegisterChannel: make(chan *clientInfo),
		clientRemoveChannel:   make(chan *websocket.Conn),
		broadcastMessage:      make(chan Message, 100), // buffer 100 messages
	}
	go h.run()

	app := fiber.New()
	app.Use(AllowUpgrade)
	app.Get("ws/direct", websocket.New(DirectMessage(h)))

	log.Info("Server started on port 3000")
	log.Fatal(app.Listen(":3000"))
}
