package repository

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type SummaryRepo struct {
	db *gorm.DB
}

func NewSummaryRepo(db *gorm.DB) *SummaryRepo {
	return &SummaryRepo{db: db}
}

func (r *SummaryRepo) Create(s *models.Summary) error {
	err := r.db.Create(s).Error
	if err != nil {
		log.Printf("[summary_repo:Create] ERROR docID=%s err=%v", s.DocumentID, err)
	} else {
		log.Printf("[summary_repo:Create] OK id=%s docID=%s", s.ID, s.DocumentID)
	}
	return err
}

func (r *SummaryRepo) FindBySubjectID(subjectID uuid.UUID) ([]models.Summary, error) {
	var summaries []models.Summary
	err := r.db.Where("subject_id = ?", subjectID).Order("created_at DESC").Find(&summaries).Error
	return summaries, err
}

func (r *SummaryRepo) FindByID(id uuid.UUID) (*models.Summary, error) {
	var s models.Summary
	err := r.db.First(&s, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SummaryRepo) FindByDocumentID(docID uuid.UUID) (*models.Summary, error) {
	var s models.Summary
	err := r.db.Where("document_id = ?", docID).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[summary_repo:FindByDocumentID] NOT FOUND docID=%s", docID)
			return nil, nil
		}
		log.Printf("[summary_repo:FindByDocumentID] ERROR docID=%s err=%v", docID, err)
		return nil, err
	}
	return &s, nil
}

func (r *SummaryRepo) UpdateContent(id uuid.UUID, short, medium, long string) error {
	return r.db.Model(&models.Summary{}).Where("id = ?", id).
		Updates(map[string]interface{}{"content_short": short, "content_medium": medium, "content_long": long}).Error
}

func (r *SummaryRepo) UpdateContentField(id uuid.UUID, field, value string) error {
	return r.db.Model(&models.Summary{}).Where("id = ?", id).Update(field, value).Error
}
