package worker

import (
	"log"
	"time"

	"github.com/google/uuid"

	"quikram-service/internal/ai"
	"quikram-service/internal/models"
	"quikram-service/internal/repository"
	"quikram-service/internal/storage"
)

type Worker struct {
	docRepo     *repository.DocumentRepo
	summaryRepo *repository.SummaryRepo
	ai          ai.LLMProvider
	storage     storage.FileStorage
}

func NewWorker(docRepo *repository.DocumentRepo, summaryRepo *repository.SummaryRepo, ai ai.LLMProvider, storage storage.FileStorage) *Worker {
	return &Worker{
		docRepo:     docRepo,
		summaryRepo: summaryRepo,
		ai:          ai,
		storage:     storage,
	}
}

func (w *Worker) ProcessDocument(documentID uuid.UUID) {
	doc, err := w.docRepo.FindByID(documentID)
	if err != nil {
		log.Printf("worker: document not found %s: %v", documentID, err)
		return
	}

	// Use already-extracted text if available (e.g. when filling in missing levels)
	content := doc.ExtractedText
	if content == "" {
		reader, err := w.storage.Download(doc.FileKey)
		if err != nil {
			log.Printf("worker: failed to download %s: %v", doc.FileKey, err)
			w.docRepo.UpdateStatus(doc.ID, "error")
			return
		}
		defer reader.Close()

		content, err = extractText(reader, doc.MimeType)
		if err != nil {
			log.Printf("worker: failed to extract text %s: %v", doc.ID, err)
			w.docRepo.UpdateStatus(doc.ID, "error")
			return
		}

		if err := w.docRepo.UpdateExtractedText(doc.ID, content); err != nil {
			log.Printf("worker: failed to save text %s: %v", doc.ID, err)
			return
		}
	}

	// Find or create summary record before AI calls
	s, err := w.summaryRepo.FindByDocumentID(doc.ID)
	if err != nil {
		log.Printf("worker: failed to check summary %s: %v", doc.ID, err)
		return
	}
	if s == nil {
		s = &models.Summary{
			DocumentID: doc.ID,
			SubjectID:  doc.SubjectID,
			UserID:     doc.UserID,
		}
		if err := w.summaryRepo.Create(s); err != nil {
			log.Printf("worker: failed to create summary %s: %v", doc.ID, err)
			return
		}
	}

	levels := []struct {
		field string
		level string
	}{
		{"content_short", "short"},
		{"content_medium", "medium"},
		{"content_long", "long"},
	}

	needsDelay := false
	for _, l := range levels {
		if s.ContentShort != "" && l.field == "content_short" {
			log.Printf("worker: %s already exists for %s, skipping", l.level, doc.ID)
			continue
		}
		if s.ContentMedium != "" && l.field == "content_medium" {
			log.Printf("worker: %s already exists for %s, skipping", l.level, doc.ID)
			continue
		}
		if s.ContentLong != "" && l.field == "content_long" {
			log.Printf("worker: %s already exists for %s, skipping", l.level, doc.ID)
			continue
		}

		if needsDelay {
			log.Printf("worker: waiting 5s before %s for %s", l.level, doc.ID)
			time.Sleep(5 * time.Second)
		}
		needsDelay = true

		result, err := w.ai.GenerateSummary(content, l.level)
		if err != nil {
			log.Printf("worker: %s summary failed %s: %v", l.level, doc.ID, err)
		} else {
			w.summaryRepo.UpdateContentField(s.ID, l.field, result)
			log.Printf("worker: %s summary saved for %s (%d chars)", l.level, doc.ID, len(result))
		}
	}

	log.Printf("worker: document %s done", doc.ID)
}

func (w *Worker) CleanupExpiredFiles() {
	docs, err := w.docRepo.FindExpired()
	if err != nil {
		log.Printf("worker: failed to find expired files: %v", err)
		return
	}

	for _, doc := range docs {
		if err := w.storage.Delete(doc.FileKey); err != nil {
			log.Printf("worker: failed to delete file %s: %v", doc.FileKey, err)
		}
		log.Printf("worker: deleted expired file %s", doc.FileKey)
	}
}

func (w *Worker) StartCleanupCron(stop chan struct{}) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.CleanupExpiredFiles()
		case <-stop:
			return
		}
	}
}
