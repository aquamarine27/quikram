package repository

import (
	"errors"

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
	return r.db.Create(s).Error
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

func (r *SummaryRepo) UpdateContent(id uuid.UUID, content, compressionLevel string) error {
	return r.db.Model(&models.Summary{}).Where("id = ?", id).
		Updates(map[string]interface{}{"content": content, "compression_level": compressionLevel}).Error
}
