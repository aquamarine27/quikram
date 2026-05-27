package ai

import (
	"quikram-service/internal/models"
)

type LLMProvider interface {
	GenerateSummary(text, compressionLevel string) (string, error)
	GenerateQuiz(summaryText string, count int, difficulty string) ([]models.Question, error)
	ChatCompletion(messages []models.ChatMessage, systemContext string) (string, error)
}
