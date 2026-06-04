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
	base := "\n\nТекст:\n"
	switch level {
	case "short":
		return "Напиши краткий конспект текста на русском языке. Выдели 2-3 смысловых блока. Каждый блок начни с короткого заголовка с новой строки. В ответе должен быть ТОЛЬКО текст конспекта. Никаких размышлений, планов, пояснений или рассуждений. Ни слова на английском." + base
	case "medium":
		return "Напиши конспект текста среднего размера на русском языке. Выдели основные темы, каждую с 2-4 предложениями. Заголовки разделов с новой строки. Между разделами пустая строка. Ключевые термины в кавычках. В ответе ТОЛЬКО конспект. Никаких размышлений и пояснений." + base
	case "long":
		return "Напиши подробный конспект текста на русском языке. Сохрани все разделы, примеры, термины, определения. Каждый раздел с заголовком на отдельной строке. Между разделами пустая строка. Ключевые термины в кавычках. В ответе ТОЛЬКО конспект. Никаких размышлений, планов или комментариев." + base
	}
	return ""
}

func GenerateSystemPrompt() string {
	return "Ты — ассистент для создания конспектов. Отвечай только готовым конспектом. Никаких рассуждений, планов, пояснений, комментариев. Никакого английского языка. Только русский. Только конспект."
}

func MaxTokensForLevel(level string) int {
	switch level {
	case "short":
		return 2048
	case "medium":
		return 4096
	case "long":
		return 8192
	}
	return 2048
}
