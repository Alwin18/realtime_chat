package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/websoket-chat/internal/repository"
)

type ApiHandler struct {
	chatMessageRepo *repository.ChatMessageRepository
}

func NewApiHandler(chatMessageRepo *repository.ChatMessageRepository) *ApiHandler {
	return &ApiHandler{
		chatMessageRepo: chatMessageRepo,
	}
}

func HelloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func (a *ApiHandler) GetChatHistory(c *fiber.Ctx) error {
	senderID := c.Query("senderId")
	receiverID := c.Query("receiverId")

	if senderID == "" || receiverID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "senderId and receiverId are required",
		})
	}

	response, err := a.chatMessageRepo.GetChatHistory(senderID, receiverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(response)
}
