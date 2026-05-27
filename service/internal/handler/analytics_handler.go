package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quikram-service/internal/repository"
)

type AnalyticsHandler struct {
	attemptRepo *repository.AttemptRepo
	subjectRepo *repository.SubjectRepo
}

func NewAnalyticsHandler(attemptRepo *repository.AttemptRepo, subjectRepo *repository.SubjectRepo) *AnalyticsHandler {
	return &AnalyticsHandler{
		attemptRepo: attemptRepo,
		subjectRepo: subjectRepo,
	}
}

func (h *AnalyticsHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	subjects, _ := h.subjectRepo.FindByUserID(userID)
	totalTests, avgScore, _ := h.attemptRepo.GetUserStats(userID)

	return c.JSON(fiber.Map{
		"subjects_count": len(subjects),
		"total_tests":    totalTests,
		"average_score":  avgScore,
	})
}

func (h *AnalyticsHandler) GetBySubject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}

	subject, err := h.subjectRepo.FindByID(subjectID)
	if err != nil || subject.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	return c.JSON(fiber.Map{
		"subject_id": subjectID,
		"message":    "Detailed analytics coming soon",
	})
}
