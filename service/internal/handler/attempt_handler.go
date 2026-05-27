package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quikram-service/internal/models"
	"quikram-service/internal/repository"
)

type AttemptHandler struct {
	attemptRepo *repository.AttemptRepo
	quizRepo    *repository.QuizRepo
}

func NewAttemptHandler(attemptRepo *repository.AttemptRepo, quizRepo *repository.QuizRepo) *AttemptHandler {
	return &AttemptHandler{attemptRepo: attemptRepo, quizRepo: quizRepo}
}

func (h *AttemptHandler) Start(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	quizID, err := uuid.Parse(c.Params("quizId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid quiz id"})
	}

	quiz, err := h.quizRepo.FindByID(quizID)
	if err != nil || quiz.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "quiz not found"})
	}

	attempt := &models.QuizAttempt{
		QuizID: quizID,
		UserID: userID,
	}
	if err := h.attemptRepo.Create(attempt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to start attempt"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"attempt_id": attempt.ID})
}

func (h *AttemptHandler) Submit(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	attemptID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid attempt id"})
	}

	attempt, err := h.attemptRepo.FindByID(attemptID)
	if err != nil || attempt.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "attempt not found"})
	}

	if attempt.CompletedAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "attempt already completed"})
	}

	var req struct {
		Answers []models.QuestionAnswer `json:"answers"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	questions, err := h.quizRepo.FindQuestionsByQuizID(attempt.QuizID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch questions"})
	}

	questionMap := make(map[uuid.UUID]models.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	correctCount := 0
	weakTopics := make(map[string]*models.WeakTopic)
	topicCount := make(map[string]int)

	for _, answer := range req.Answers {
		q, exists := questionMap[answer.QuestionID]
		if !exists {
			continue
		}

		correct := true
		for _, optIdx := range answer.SelectedOptionIDs {
			if optIdx >= len(q.Options) || !q.Options[optIdx].IsCorrect {
				correct = false
				break
			}
		}

		allCorrectChoices := true
		for i, opt := range q.Options {
			if opt.IsCorrect {
				selected := false
				for _, idx := range answer.SelectedOptionIDs {
					if idx == i {
						selected = true
						break
					}
				}
				if !selected {
					allCorrectChoices = false
					break
				}
			}
		}
		correct = correct && allCorrectChoices

		answer.IsCorrect = correct
		if correct {
			correctCount++
		}

		if q.Topic != "" {
			topicCount[q.Topic]++
			if _, ok := weakTopics[q.Topic]; !ok {
				weakTopics[q.Topic] = &models.WeakTopic{Topic: q.Topic}
			}
			wt := weakTopics[q.Topic]
			wt.Score += map[bool]float64{true: 100, false: 0}[correct]
		}
	}

	score := 0.0
	if len(req.Answers) > 0 {
		score = float64(correctCount) / float64(len(req.Answers)) * 100
	}

	var weakList []models.WeakTopic
	for _, wt := range weakTopics {
		if cnt, ok := topicCount[wt.Topic]; ok {
			wt.Score /= float64(cnt)
		}
		if wt.Score < 60 {
			weakList = append(weakList, *wt)
		}
	}

	now := time.Now()
	attempt.Answers = req.Answers
	attempt.Score = score
	attempt.WeakTopics = weakList
	attempt.CompletedAt = &now

	if err := h.attemptRepo.Update(attempt); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save results"})
	}

	return c.JSON(fiber.Map{
		"score":       score,
		"correct":     correctCount,
		"total":       len(req.Answers),
		"weak_topics": weakList,
	})
}

func (h *AttemptHandler) History(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	quizID, err := uuid.Parse(c.Params("quizId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid quiz id"})
	}

	quiz, err := h.quizRepo.FindByID(quizID)
	if err != nil || quiz.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "quiz not found"})
	}

	attempts, err := h.attemptRepo.FindByQuizID(quizID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch attempts"})
	}

	return c.JSON(attempts)
}

func (h *AttemptHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	attemptID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid attempt id"})
	}

	attempt, err := h.attemptRepo.FindByID(attemptID)
	if err != nil || attempt.UserID != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "attempt not found"})
	}

	return c.JSON(attempt)
}
