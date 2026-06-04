package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"quikram-service/internal/models"
)

type OpenAIProvider struct {
	apiKey   string
	endpoint string
	model    string
	client   *http.Client
}

func NewOpenAIProvider(apiKey, endpoint, model string) *OpenAIProvider {
	if model == "" {
		if strings.Contains(endpoint, "openrouter") {
			model = "nvidia/nemotron-3-super-120b-a12b:free"
		} else if strings.Contains(endpoint, "groq") {
			model = "llama-3.1-70b-versatile"
		} else {
			model = "gpt-4o-mini"
		}
	}

	return &OpenAIProvider{
		apiKey:   apiKey,
		endpoint: strings.TrimRight(endpoint, "/"),
		model:    model,
		client:   &http.Client{Timeout: 600 * time.Second},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *OpenAIProvider) chatCompletion(messages []chatMessage, maxTokens int) (string, error) {
	body := chatRequest{
		Model:       p.model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: 0.7,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", p.endpoint+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("api request failed: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("api returned %d: %s", resp.StatusCode, string(respData))
		if resp.StatusCode == 402 {
			errMsg = "API insufficient credits — check your OpenRouter billing: https://openrouter.ai/settings/credits"
		}
		return "", fmt.Errorf(errMsg)
	}

	var result chatResponse
	if err := json.Unmarshal(respData, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("api error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

func (p *OpenAIProvider) GenerateSummary(text, level string) (string, error) {
	prompt := GeneratePrompt(level) + text
	result, err := p.chatCompletion([]chatMessage{
		{Role: "system", Content: GenerateSystemPrompt()},
		{Role: "user", Content: prompt},
	}, MaxTokensForLevel(level))
	if err != nil {
		return "", err
	}
	log.Printf("[openai:GenerateSummary] level=%s length=%d", level, len(result))
	return result, nil
}

func (p *OpenAIProvider) GenerateQuiz(summaryText string, count int, difficulty string) ([]models.Question, error) {
	difficultyPrompt := "средней"
	if difficulty == "easy" {
		difficultyPrompt = "лёгкой"
	} else if difficulty == "hard" {
		difficultyPrompt = "сложной"
	}

	prompt := fmt.Sprintf(`На основе текста ниже составь тест из %d вопросов %s сложности. Формат ответа — строгий JSON-массив:

[
  {
    "question_text": "текст вопроса",
    "question_type": "single_choice",
    "options": [
      {"text": "вариант 1", "is_correct": false},
      {"text": "вариант 2", "is_correct": true},
      {"text": "вариант 3", "is_correct": false},
      {"text": "вариант 4", "is_correct": false}
    ],
    "explanation": "почему этот ответ правильный",
    "topic": "тема вопроса"
  }
]

Текст:
%s

JSON:`, count, difficultyPrompt, summaryText)

	resp, err := p.chatCompletion([]chatMessage{
		{Role: "user", Content: prompt},
	}, 4096)
	if err != nil {
		return nil, err
	}

	resp = cleanJSON(resp)

	var questions []models.Question
	if err := json.Unmarshal([]byte(resp), &questions); err != nil {
		return nil, fmt.Errorf("failed to parse quiz JSON: %w\nraw: %s", err, resp)
	}

	for i := range questions {
		questions[i].OrderIndex = i
	}

	return questions, nil
}

func (p *OpenAIProvider) ChatCompletion(messages []models.ChatMessage, systemContext string) (string, error) {
	chatMessages := []chatMessage{}
	if systemContext != "" {
		chatMessages = append(chatMessages, chatMessage{Role: "system", Content: systemContext})
	}
	for _, m := range messages {
		chatMessages = append(chatMessages, chatMessage{Role: m.Role, Content: m.Content})
	}

	return p.chatCompletion(chatMessages, 2048)
}

func cleanJSON(raw string) string {
	start := strings.Index(raw, "[")
	if start == -1 {
		return raw
	}
	end := strings.LastIndex(raw, "]")
	if end == -1 || end < start {
		return raw
	}
	return raw[start : end+1]
}
