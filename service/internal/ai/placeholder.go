package ai

import (
	"fmt"
	"math/rand"

	"quikram-service/internal/models"
)

type PlaceholderProvider struct{}

func NewPlaceholderProvider() *PlaceholderProvider {
	return &PlaceholderProvider{}
}

func (p *PlaceholderProvider) GenerateSummary(text, level string) (string, error) {
	info := fmt.Sprintf("Файл содержит %d символов", len(text))
	switch level {
	case "short":
		return "Короткий конспект\n\nЭто placeholder. Подключи API-ключ для генерации.\n" + info, nil
	case "medium":
		return "Средний конспект\n\nЭто placeholder-конспект. Текст будет генерироваться ИИ после подключения API.\n\nОсновные темы\nОписание темы 1\nОписание темы 2\nОписание темы 3\n\n" + info, nil
	case "long":
		return "Длинный конспект\n\nЭто placeholder-конспект. Текст будет генерироваться ИИ после подключения API.\n\nОсновные темы\nОписание темы 1 с деталями\nОписание темы 2 с деталями\nОписание темы 3 с деталями\n\nКлючевые определения\nТермин\n" + info, nil
	}
	return "", fmt.Errorf("unknown level: %s", level)
}

func (p *PlaceholderProvider) GenerateQuiz(summaryText string, count int, difficulty string) ([]models.Question, error) {
	questions := make([]models.Question, count)
	for i := 0; i < count; i++ {
		correctIdx := rand.Intn(4)
		options := make([]models.QuestionOption, 4)
		for j := 0; j < 4; j++ {
			options[j] = models.QuestionOption{
				Text:      fmt.Sprintf("Вариант ответа %d", j+1),
				IsCorrect: j == correctIdx,
			}
		}

		questions[i] = models.Question{
			QuestionText: fmt.Sprintf("Placeholder вопрос %d? Выберите правильный ответ.", i+1),
			QuestionType: "single_choice",
			Options:      options,
			Explanation:  fmt.Sprintf("Это объяснение правильного ответа на вопрос %d.", i+1),
			Topic:        fmt.Sprintf("Тема %d", (i%3)+1),
			OrderIndex:   i,
		}
	}
	return questions, nil
}

func (p *PlaceholderProvider) ChatCompletion(messages []models.ChatMessage, systemContext string) (string, error) {
	return "Это placeholder-ответ от AI-ассистента. Подключите API-ключ для реальных ответов.", nil
}
