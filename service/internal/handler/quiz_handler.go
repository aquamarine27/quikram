package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quikram-service/internal/ai"
	"quikram-service/internal/models"
	"quikram-service/internal/repository"
	"quikram-service/internal/service"
	"quikram-service/internal/validator"
)

type QuizHandler struct {
	quizRepo       *repository.QuizRepo
	summaryRepo    *repository.SummaryRepo
	ai             ai.LLMProvider
	subjectService *service.SubjectService
}

func NewQuizHandler(quizRepo *repository.QuizRepo, summaryRepo *repository.SummaryRepo, ai ai.LLMProvider, subjectService *service.SubjectService) *QuizHandler {
	return &QuizHandler{
		quizRepo:       quizRepo,
		summaryRepo:    summaryRepo,
		ai:             ai,
		subjectService: subjectService,
	}
}

func (h *QuizHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	var req struct {
		SummaryID      uuid.UUID `json:"summary_id"`
		QuestionsCount int       `json:"questions_count"`
		Difficulty     string    `json:"difficulty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}
	if err := validator.Difficulty(req.Difficulty); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.QuestionsCount < 5 || req.QuestionsCount > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "questions count must be between 5 and 30"})
	}

	summary, err := h.summaryRepo.FindByID(req.SummaryID)
	if err != nil || summary.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "summary not found"})
	}

	questions, err := h.ai.GenerateQuiz(summary.Content, req.QuestionsCount, req.Difficulty)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate quiz"})
	}

	quiz := &models.Quiz{
		SummaryID:      req.SummaryID,
		SubjectID:      subjectID,
		UserID:         userID,
		Title:          "Тест по конспекту",
		QuestionsCount: req.QuestionsCount,
		Difficulty:     req.Difficulty,
	}
	if err := h.quizRepo.Create(quiz); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save quiz"})
	}

	for i := range questions {
		questions[i].QuizID = quiz.ID
	}
	if err := h.quizRepo.CreateQuestions(questions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save questions"})
	}

	return c.Status(fiber.StatusCreated).JSON(quiz)
}

func (h *QuizHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	quizzes, err := h.quizRepo.FindBySubjectID(subjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch quizzes"})
	}

	return c.JSON(quizzes)
}

func (h *QuizHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	quizID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid quiz id"})
	}

	quiz, err := h.quizRepo.FindByID(quizID)
	if err != nil || quiz.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "quiz not found"})
	}

	questions, err := h.quizRepo.FindQuestionsByQuizID(quizID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch questions"})
	}

	for i := range questions {
		for j := range questions[i].Options {
			questions[i].Options[j].IsCorrect = false
		}
	}

	return c.JSON(fiber.Map{
		"quiz":      quiz,
		"questions": questions,
	})
}

func (h *QuizHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	quizID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid quiz id"})
	}

	quiz, err := h.quizRepo.FindByID(quizID)
	if err != nil || quiz.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "quiz not found"})
	}

	if err := h.quizRepo.Delete(quizID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete quiz"})
	}

	return c.JSON(fiber.Map{"message": "quiz deleted"})
}
