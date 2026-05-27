package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quikram-service/internal/repository"
	"quikram-service/internal/service"
	"quikram-service/internal/validator"
)

type SummaryHandler struct {
	summaryRepo    *repository.SummaryRepo
	subjectService *service.SubjectService
}

func NewSummaryHandler(summaryRepo *repository.SummaryRepo, subjectService *service.SubjectService) *SummaryHandler {
	return &SummaryHandler{summaryRepo: summaryRepo, subjectService: subjectService}
}

func (h *SummaryHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	summaries, err := h.summaryRepo.FindBySubjectID(subjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch summaries"})
	}

	return c.JSON(summaries)
}

func (h *SummaryHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	summaryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid summary id"})
	}

	summary, err := h.summaryRepo.FindByID(summaryID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "summary not found"})
	}

	if summary.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	return c.JSON(summary)
}

func (h *SummaryHandler) Regenerate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	summaryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid summary id"})
	}

	var req struct {
		CompressionLevel string `json:"compression_level"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.CompressionLevel == "" {
		req.CompressionLevel = "medium"
	}
	if err := validator.CompressionLevel(req.CompressionLevel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	summary, err := h.summaryRepo.FindByID(summaryID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "summary not found"})
	}

	if summary.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	return c.JSON(fiber.Map{"message": "regeneration started"})
}
