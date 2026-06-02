package handler

import "github.com/gofiber/fiber/v2"

type Review struct {
	Text   string `json:"text"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Badge  string `json:"badge"`
}

type ReviewHandler struct{}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

func (h *ReviewHandler) List(c *fiber.Ctx) error {
	lang := c.Query("lang", "ru")

	names := []string{
		"Алексей", "Мария", "Дмитрий", "Джеймс",
		"Екатерина", "Сара", "Иван", "Майкл",
		"Анна", "Лена", "Оливер", "Никита",
	}
	authors := []string{
		"МИРЭА, 3 курс", "МГИМО, 4 курс", "МГУ, 1 курс", "MIT, 2 курс",
		"ВШЭ, 4 курс", "Oxford, 3 курс", "МФТИ, 2 курс", "Stanford, 1 курс",
		"СПбГУ, 3 курс", "TU Berlin, Магистратура", "Cambridge, Выпускник", "ДВФУ, 2 курс",
	}
	texts := []string{
		"Quikram помог мне подготовиться к экзаменам за считанные недели. Удобный интерфейс, понятные конспекты и тесты — всё что нужно для эффективной подготовки.",
		"Благодаря Quikram я систематизировал все свои материалы и смог сосредоточиться на самых важных темах.",
		"AI-помощник и тестирование с аналитикой слабых мест — это именно то, чего не хватало в подготовке. Очень рекомендую!",
		"Автоматическое создание тестов по моим лекциям сэкономило часы подготовки. Просто загружаю файлы — получаю готовые вопросы.",
		"Всё в одном месте: конспекты, тесты, статистика. Не нужно прыгать между десятком приложений.",
		"AI-помощник объяснил сложную тему простыми словами за 5 минут. Преподаватели так не умеют.",
		"Тесты с адаптивным уровнем сложности заставляют думать. Не просто зубрёжка, а реальное понимание.",
		"После месяца занятий средний балл вырос на 40%. Результат говорит сам за себя.",
		"Наконец-то сервис, где можно загрузить свои файлы и получить готовые конспекты. Работает идеально.",
		"Статистика по темам показывает именно те места, где у меня пробелы. Очень помогает в подготовке.",
		"Использую Quikram для подготовки к поступлению в магистратуру. Чётко, структурированно, эффективно.",
		"Здорово, что можно повторять материал с телефона в любое время. Синхронизация с Telegram ботом — топ.",
	}

	if lang == "en" {
		names = []string{
			"Alexey", "Maria", "Dmitry", "James",
			"Catherine", "Sarah", "Ivan", "Michael",
			"Anna", "Lena", "Oliver", "Nikita",
		}
		authors = []string{
			"MIREA, 3rd year", "MGIMO, 4th year", "MSU, 1st year", "MIT, 2nd year",
			"HSE, 4th year", "Oxford, 3rd year", "MIPT, 2nd year", "Stanford, 1st year",
			"SPbSU, 3rd year", "TU Berlin, Master", "Cambridge, Graduate", "FEFU, 2nd year",
		}
		texts = []string{
			"Quikram helped me prepare for exams in just a few weeks. Convenient interface, clear notes and tests — everything I needed.",
			"Thanks to Quikram I organized all my materials and focused on the most important topics.",
			"The AI assistant with weak spot analytics is exactly what I was missing. Highly recommend!",
			"Auto-generated tests from my lectures saved hours of prep time. Just upload files — get ready questions.",
			"Everything in one place: notes, tests, stats. No more jumping between a dozen apps.",
			"The AI explained a complex topic in simple words in 5 minutes. Teachers can't do that.",
			"Adaptive difficulty tests make you think. Real understanding, not just memorization.",
			"After a month of using Quikram my average score went up 40%. Results speak for themselves.",
			"Finally a service where I can upload my own files and get structured notes. Works perfectly.",
			"Topic analytics pinpoint exactly where my weak spots are. Incredibly helpful for focused prep.",
			"I use Quikram to prepare for my master's entrance exams. Clear, structured, effective.",
			"Love that I can review material on my phone anytime. Telegram bot sync is awesome.",
		}
	}

	badges := []string{"#8b5cf6", "#f59e0b", "#10b981", "#3b82f6", "#ef4444", "#ec4899"}

	reviews := make([]Review, 12)
	for i := 0; i < 12; i++ {
		reviews[i] = Review{
			Text:   texts[i],
			Name:   names[i],
			Author: authors[i],
			Badge:  badges[i%6],
		}
	}

	return c.JSON(reviews)
}
