package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"quikram-service/internal/ai"
	"quikram-service/internal/config"
	"quikram-service/internal/db"
	"quikram-service/internal/handler"
	"quikram-service/internal/middleware"
	"quikram-service/internal/repository"
	"quikram-service/internal/service"
	"quikram-service/internal/storage"
	"quikram-service/internal/worker"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("database migrations completed")

	userRepo := repository.NewUserRepo(database)
	subjectRepo := repository.NewSubjectRepo(database)
	documentRepo := repository.NewDocumentRepo(database)
	summaryRepo := repository.NewSummaryRepo(database)
	quizRepo := repository.NewQuizRepo(database)
	attemptRepo := repository.NewAttemptRepo(database)
	tokenRepo := repository.NewTokenRepo(database)
	chatRepo := repository.NewChatRepo(database)

	llmProvider := ai.NewPlaceholderProvider()

	fileStorage := storage.NewLocalStorage("./uploads")

	workerPool := worker.NewWorker(documentRepo, summaryRepo, llmProvider, fileStorage)
	stopCron := make(chan struct{})
	go workerPool.StartCleanupCron(stopCron)

	authService := service.NewAuthService(userRepo, tokenRepo, cfg)
	subjectService := service.NewSubjectService(subjectRepo)

	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userRepo)
	subjectHandler := handler.NewSubjectHandler(subjectService)
	documentHandler := handler.NewDocumentHandler(documentRepo, subjectService, fileStorage, workerPool)
	summaryHandler := handler.NewSummaryHandler(summaryRepo, subjectService)
	quizHandler := handler.NewQuizHandler(quizRepo, summaryRepo, llmProvider, subjectService)
	attemptHandler := handler.NewAttemptHandler(attemptRepo, quizRepo)
	analyticsHandler := handler.NewAnalyticsHandler(attemptRepo, subjectRepo)
	planHandler := handler.NewPlanHandler()
	chatHandler := handler.NewChatHandler(chatRepo, llmProvider)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal error",
			})
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.CORS(cfg.FrontendURL))
	app.Use(middleware.RateLimit(100, 60))

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	authMw := middleware.AuthRequired(cfg.JWTSecret)

	users := api.Group("/users", authMw)
	users.Get("/me", userHandler.GetMe)
	users.Patch("/me", userHandler.UpdateMe)
	users.Post("/me/change-password", userHandler.ChangePassword)
	users.Post("/me/change-plan", userHandler.ChangePlan)

	subjects := api.Group("/subjects", authMw)
	subjects.Get("/", subjectHandler.List)
	subjects.Post("/", subjectHandler.Create)
	subjects.Get("/:id", subjectHandler.GetByID)
	subjects.Patch("/:id", subjectHandler.Update)
	subjects.Delete("/:id", subjectHandler.Delete)

	documents := api.Group("/subjects/:subjectId/documents", authMw)
	documents.Post("/", documentHandler.Upload)
	documents.Get("/", documentHandler.List)
	documents.Get("/:id", documentHandler.GetByID)
	documents.Delete("/:id", documentHandler.Delete)

	summaries := api.Group("/subjects/:subjectId/summaries", authMw)
	summaries.Get("/", summaryHandler.List)
	summaries.Get("/:id", summaryHandler.GetByID)
	summaries.Post("/:id/regenerate", summaryHandler.Regenerate)

	quizzes := api.Group("/subjects/:subjectId/quizzes", authMw)
	quizzes.Post("/", quizHandler.Create)
	quizzes.Get("/", quizHandler.List)
	quizzes.Get("/:id", quizHandler.GetByID)
	quizzes.Delete("/:id", quizHandler.Delete)

	attempts := api.Group("/quizzes/:quizId/attempts", authMw)
	attempts.Post("/", attemptHandler.Start)
	attempts.Post("/:id/submit", attemptHandler.Submit)
	attempts.Get("/", attemptHandler.History)
	attempts.Get("/:id", attemptHandler.GetByID)

	chat := api.Group("/chat", authMw)
	chat.Post("/", chatHandler.Send)
	chat.Get("/", chatHandler.History)
	chat.Delete("/", chatHandler.Clear)

	api.Get("/plans", planHandler.List)

	analytics := api.Group("/analytics", authMw)
	analytics.Get("/me", analyticsHandler.GetMe)
	analytics.Get("/subjects/:id", analyticsHandler.GetBySubject)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		close(stopCron)
		app.Shutdown()
	}()

	log.Printf("server starting on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
