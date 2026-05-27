package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quikram-service/internal/ai"
	"quikram-service/internal/repository"
)

type ChatHandler struct {
	chatRepo *repository.ChatRepo
	ai       ai.LLMProvider
}

func NewChatHandler(chatRepo *repository.ChatRepo, ai ai.LLMProvider) *ChatHandler {
	return &ChatHandler{chatRepo: chatRepo, ai: ai}
}

func (h *ChatHandler) Send(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req struct {
		Message string `json:"message"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userMsg, err := h.chatRepo.Save(userID, "user", req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save message"})
	}

	history, _ := h.chatRepo.FindLastByUserID(userID, 10)

	systemContext := "Ты умный AI-ассистент. Отвечай полезно и по делу."

	aiResponse, err := h.ai.ChatCompletion(history, systemContext)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "AI service error"})
	}

	aiMsg, err := h.chatRepo.Save(userID, "assistant", aiResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save response"})
	}

	return c.JSON(fiber.Map{
		"user":      userMsg,
		"assistant": aiMsg,
	})
}

func (h *ChatHandler) History(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	messages, err := h.chatRepo.FindByUserID(userID, 50)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch messages"})
	}

	return c.JSON(messages)
}

func (h *ChatHandler) Clear(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	if err := h.chatRepo.DeleteByUserID(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to clear chat"})
	}

	return c.JSON(fiber.Map{"message": "chat history cleared"})
}
