package ai

import (
	"quikram-service/internal/models"
)

type LLMProvider interface {
	GenerateSummary(text, level string) (string, error)
	GenerateQuiz(summaryText string, count int, difficulty string) ([]models.Question, error)
	ChatCompletion(messages []models.ChatMessage, systemContext string) (string, error)
}

func GeneratePrompt(level string) string {
	switch level {
	case "short":
		return "Напиши краткий конспект текста. Основные мысли, хватит для быстрого понимания. Раздели на 2-3 смысловых блока, каждый с коротким заголовком на отдельной строке. Без вступлений, без воды."
	case "medium":
		return "Напиши средний конспект текста. Выдели основные темы, каждую с 2-4 предложениями. Заголовки разделов пиши с новой строки (без символа #). Между разделами пустая строка. Ключевые термины в кавычках."
	case "long":
		return "Напиши подробный конспект текста. Сохрани все разделы, примеры, термины, определения. Каждый раздел начинай с заголовка на отдельной строке (без символа #). Используй пустые строки между разделами. Ключевые термины в кавычках. Детально, со всеми важными деталями."
	}
	return ""
}

func MaxTokensForLevel(level string) int {
	switch level {
	case "short":
		return 1024
	case "medium":
		return 3072
	case "long":
		return 4096
	}
	return 1024
}
