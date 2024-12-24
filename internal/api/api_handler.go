package api

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/websoket-chat/internal/repository"
)

type ApiHandler struct {
	chatMessageRepo *repository.ChatMessageRepository
	contactRepo     *repository.ContactRepository
}

func NewApiHandler(
	chatMessageRepo *repository.ChatMessageRepository,
	contactRepo *repository.ContactRepository,
) *ApiHandler {
	return &ApiHandler{
		chatMessageRepo: chatMessageRepo,
		contactRepo:     contactRepo,
	}
}

func HelloHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

func (a *ApiHandler) GetChatHistory(c *fiber.Ctx) error {
	if c.Query("senderId", "") == "" || c.Query("receiverId", "") == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messages": "senderId and receiverId are required",
			"code":     fiber.StatusBadRequest,
		})
	}

	response, err := a.chatMessageRepo.GetChatHistory(c.Query("senderId", ""), c.Query("receiverId", ""))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"data":    response,
		"status":  true,
		"message": "success get history chat",
		"code":    fiber.StatusOK,
	})
}

func (a *ApiHandler) GetContactByCakupan(c *fiber.Ctx) error {
	if c.Query("gedung_id", "") == "" || c.Query("role", "") == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"messages": "gedung_id and role are required",
			"code":     fiber.StatusBadRequest,
		})
	}

	gedungId, err := strconv.ParseInt(c.Query("gedung_id", ""), 10, 64)
	if err != nil {
		return c.JSON(fiber.Map{
			"messages": errors.New("gedung_id must be integer"),
			"code":     fiber.StatusInternalServerError,
		})
	}

	response, err := a.contactRepo.GetContactByCakupan(gedungId, c.Query("role", ""))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusInternalServerError,
		})
	}

	return c.JSON(fiber.Map{
		"data":    response,
		"status":  true,
		"message": "success get list contacts",
		"code":    fiber.StatusOK,
	})
}
