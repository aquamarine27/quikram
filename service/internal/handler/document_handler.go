package handler

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quikram-service/internal/models"
	"quikram-service/internal/repository"
	"quikram-service/internal/service"
	"quikram-service/internal/storage"
	"quikram-service/internal/validator"
	"quikram-service/internal/worker"
)

type DocumentHandler struct {
	docRepo        *repository.DocumentRepo
	subjectService *service.SubjectService
	userRepo       *repository.UserRepo
	storage        storage.FileStorage
	worker         *worker.Worker
}

func NewDocumentHandler(docRepo *repository.DocumentRepo, subjectService *service.SubjectService, userRepo *repository.UserRepo, storage storage.FileStorage, worker *worker.Worker) *DocumentHandler {
	return &DocumentHandler{
		docRepo:        docRepo,
		subjectService: subjectService,
		userRepo:       userRepo,
		storage:        storage,
		worker:         worker,
	}
}

func (h *DocumentHandler) Upload(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not found"})
	}

	if user.Plan == "free" && user.UploadsThisMonth >= 15 {
		return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
			"error": "Free plan limit reached: max 15 uploads per month. Upgrade to Pro or ProAI for unlimited uploads.",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required"})
	}

	if file.Size > 20*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file too large: max 20 MB"})
	}

	mimeType := file.Header.Get("Content-Type")
	if err := validator.MimeType(mimeType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read file"})
	}
	defer f.Close()

	fileKey := fmt.Sprintf("%s/%s/%s", userID, subjectID, uuid.New().String()+filepath.Ext(file.Filename))

	if err := h.storage.Upload(fileKey, f, mimeType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to store file"})
	}

	deleteAt := time.Now().Add(24 * time.Hour)

	doc := &models.Document{
		SubjectID: subjectID,
		UserID:    userID,
		Filename:  file.Filename,
		FileKey:   fileKey,
		FileSize:  int(file.Size),
		MimeType:  mimeType,
		Status:    "uploaded",
		DeleteAt:  &deleteAt,
	}

	if err := h.docRepo.Create(doc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save document"})
	}

	if err := h.userRepo.IncrementUploads(userID); err != nil {
		log.Printf("failed to increment uploads for user %s: %v", userID, err)
	}

	go h.worker.ProcessDocument(doc.ID)

	return c.Status(fiber.StatusCreated).JSON(doc)
}

func (h *DocumentHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	docs, err := h.docRepo.FindBySubjectID(subjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch documents"})
	}

	return c.JSON(docs)
}

func (h *DocumentHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	docID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid document id"})
	}

	doc, err := h.docRepo.FindByID(docID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}

	if doc.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	return c.JSON(doc)
}

func (h *DocumentHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	docID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid document id"})
	}

	doc, err := h.docRepo.FindByID(docID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}

	if doc.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	h.storage.Delete(doc.FileKey)
	h.docRepo.Delete(docID)

	return c.JSON(fiber.Map{"message": "document deleted"})
}
