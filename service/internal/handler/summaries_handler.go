package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"quikram-service/internal/repository"
	"quikram-service/internal/service"
	"quikram-service/internal/validator"
	"quikram-service/internal/worker"
)

type SummaryHandler struct {
	summaryRepo    *repository.SummaryRepo
	subjectService *service.SubjectService
	userRepo       *repository.UserRepo
	docRepo        *repository.DocumentRepo
	worker         *worker.Worker
}

func NewSummaryHandler(summaryRepo *repository.SummaryRepo, subjectService *service.SubjectService, userRepo *repository.UserRepo, docRepo *repository.DocumentRepo, worker *worker.Worker) *SummaryHandler {
	return &SummaryHandler{summaryRepo: summaryRepo, subjectService: subjectService, userRepo: userRepo, docRepo: docRepo, worker: worker}
}

func (h *SummaryHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	summaries, err := h.summaryRepo.FindBySubjectID(subjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch summaries"})
	}

	return c.JSON(summaries)
}

func (h *SummaryHandler) GetByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	summaryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid summary id"})
	}

	summary, err := h.summaryRepo.FindByID(summaryID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "summary not found"})
	}

	if summary.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	return c.JSON(summary)
}

func (h *SummaryHandler) GetByDocument(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}
	docID, err := uuid.Parse(c.Params("docId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid document id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	doc, err := h.docRepo.FindByID(docID)
	if err != nil || doc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}
	if doc.UserID != userID || doc.SubjectID != subjectID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	log.Printf("[summary:GetByDocument] docID=%s userID=%s", docID, userID)
	summary, err := h.summaryRepo.FindByDocumentID(docID)
	if err != nil {
		log.Printf("[summary:GetByDocument] FindByDocumentID error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch summary"})
	}
	if summary == nil {
		log.Printf("[summary:GetByDocument] summary NOT FOUND for docID=%s", docID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "summary not found"})
	}

	log.Printf("[summary:GetByDocument] summary FOUND id=%s docID=%s", summary.ID, summary.DocumentID)
	return c.JSON(summary)
}

func (h *SummaryHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	subjectID, err := uuid.Parse(c.Params("subjectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid subject id"})
	}
	docID, err := uuid.Parse(c.Params("docId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid document id"})
	}

	if _, err := h.subjectService.GetByID(subjectID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subject not found"})
	}

	doc, err := h.docRepo.FindByID(docID)
	if err != nil || doc == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}
	if doc.UserID != userID || doc.SubjectID != subjectID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	existing, _ := h.summaryRepo.FindByDocumentID(docID)
	if existing != nil && (existing.ContentShort != "" || existing.ContentMedium != "" || existing.ContentLong != "") {
		return c.JSON(existing)
	}
	if existing != nil {
		// summary exists but empty — still being generated
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "processing", "message": "summary is being generated"})
	}

	if doc.Status == "error" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "file upload failed, re-upload the document"})
	}

	go h.worker.ProcessDocument(doc.ID)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "processing", "message": "summary generation started"})
}

func (h *SummaryHandler) Regenerate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	summaryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid summary id"})
	}

	var req struct {
		CompressionLevel string `json:"compression_level"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.CompressionLevel == "" {
		req.CompressionLevel = "medium"
	}
	if err := validator.CompressionLevel(req.CompressionLevel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	summary, err := h.summaryRepo.FindByID(summaryID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "summary not found"})
	}

	if summary.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "access denied"})
	}

	return c.JSON(fiber.Map{"message": "regeneration started"})
}
