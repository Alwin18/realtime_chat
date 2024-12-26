package api

import (
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
	var body GetChatHistoryRequest
	if err := c.QueryParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NewResponseWithoutData(
			"failed request: "+err.Error(),
			fiber.StatusBadRequest,
			false,
		))
	}

	if body.SenderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(NewResponseWithoutData(
			"sender_id is required",
			fiber.StatusBadRequest,
			false,
		))
	}

	if body.ReceiverID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(NewResponseWithoutData(
			"receiver_id is required",
			fiber.StatusBadRequest,
			false,
		))
	}

	response, err := a.chatMessageRepo.GetChatHistory(body.SenderID, body.ReceiverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewResponseWithoutData(
			"failed get list chats: "+err.Error(),
			fiber.StatusInternalServerError,
			false,
		))
	}

	return c.JSON(NewBaseResponse(
		true,
		"success get list contacts",
		fiber.StatusOK,
		response,
	))
}

func (a *ApiHandler) GetContactByCakupan(c *fiber.Ctx) error {
	var body GetContactByCakupanRequest
	if err := c.QueryParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(NewResponseWithoutData(
			"failed request: "+err.Error(),
			fiber.StatusBadRequest,
			false,
		))
	}

	if body.GedungID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(NewResponseWithoutData(
			"gedung_id is required",
			fiber.StatusBadRequest,
			false,
		))
	}

	if body.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(NewResponseWithoutData(
			"role is required",
			fiber.StatusBadRequest,
			false,
		))
	}

	contacts, err := a.contactRepo.GetContactByCakupan(body.GedungID, body.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NewResponseWithoutData(
			"failed get list contacts: "+err.Error(),
			fiber.StatusInternalServerError,
			false,
		))
	}

	users := make([]GetContactByCakupanResponse, 0)
	for _, v := range contacts {
		users = append(users, GetContactByCakupanResponse{
			ID:          v.ID,
			Name:        v.Name,
			Email:       v.Email,
			PhoneNumber: v.PhoneNumber,
			AvatarURL:   v.AvatarURL,
			KotaID:      v.KotaID,
			GedungID:    v.GedungID,
			IsOnline:    v.IsOnline,
		})
	}

	return c.JSON(NewBaseResponse(
		true,
		"success get list contacts",
		fiber.StatusOK,
		users,
	))
}
