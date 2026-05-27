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

	w.docRepo.UpdateStatus(doc.ID, "processing")

	reader, err := w.storage.Download(doc.FileKey)
	if err != nil {
		log.Printf("worker: failed to download %s: %v", doc.FileKey, err)
		w.docRepo.UpdateStatus(doc.ID, "error")
		return
	}
	defer reader.Close()

	content, err := extractText(reader, doc.MimeType)
	if err != nil {
		log.Printf("worker: failed to extract text %s: %v", doc.ID, err)
		w.docRepo.UpdateStatus(doc.ID, "error")
		return
	}

	if err := w.docRepo.UpdateExtractedText(doc.ID, content); err != nil {
		log.Printf("worker: failed to save text %s: %v", doc.ID, err)
		return
	}

	summary, err := w.ai.GenerateSummary(content, "medium")
	if err != nil {
		log.Printf("worker: failed to generate summary %s: %v", doc.ID, err)
		w.docRepo.UpdateStatus(doc.ID, "error")
		return
	}

	s := &models.Summary{
		DocumentID:       doc.ID,
		SubjectID:        doc.SubjectID,
		UserID:           doc.UserID,
		Content:          summary,
		CompressionLevel: "medium",
	}
	if err := w.summaryRepo.Create(s); err != nil {
		log.Printf("worker: failed to save summary %s: %v", doc.ID, err)
		return
	}

	w.docRepo.UpdateStatus(doc.ID, "ready")
	log.Printf("worker: document %s processed successfully", doc.ID)
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
