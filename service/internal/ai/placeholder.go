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

func (p *PlaceholderProvider) GenerateSummary(text, compressionLevel string) (string, error) {
	prefix := ""
	switch compressionLevel {
	case "short":
		prefix = "## Краткий конспект (short)\n\n"
	case "medium":
		prefix = "## Конспект (medium)\n\n"
	case "detailed":
		prefix = "## Подробный конспект (detailed)\n\n"
	}

	return prefix + "Это placeholder-конспект. Текст будет генерироваться ИИ после подключения API.\n\n" +
		"### Основные темы:\n- Тема 1: описание\n- Тема 2: описание\n- Тема 3: описание\n\n" +
		"### Ключевые определения:\n- **Определение 1**: значение\n- **Определение 2**: значение\n\n" +
		"_Файл содержит " + fmt.Sprintf("%d", len(text)/10) + " символов входного текста._\n", nil
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
