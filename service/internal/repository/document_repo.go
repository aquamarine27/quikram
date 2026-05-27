package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type DocumentRepo struct {
	db *gorm.DB
}

func NewDocumentRepo(db *gorm.DB) *DocumentRepo {
	return &DocumentRepo{db: db}
}

func (r *DocumentRepo) Create(d *models.Document) error {
	return r.db.Create(d).Error
}

func (r *DocumentRepo) FindBySubjectID(subjectID uuid.UUID) ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("subject_id = ?", subjectID).Order("created_at DESC").Find(&docs).Error
	return docs, err
}

func (r *DocumentRepo) FindByID(id uuid.UUID) (*models.Document, error) {
	var d models.Document
	err := r.db.First(&d, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *DocumentRepo) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.Document{}).Where("id = ?", id).Update("status", status).Error
}

func (r *DocumentRepo) UpdateExtractedText(id uuid.UUID, text string) error {
	return r.db.Model(&models.Document{}).Where("id = ?", id).
		Updates(map[string]interface{}{"extracted_text": text, "status": "ready"}).Error
}

func (r *DocumentRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Document{}, "id = ?", id).Error
}

func (r *DocumentRepo) FindExpired() ([]models.Document, error) {
	var docs []models.Document
	err := r.db.Where("delete_at IS NOT NULL AND delete_at <= ?", time.Now()).
		Select("id, file_key").Find(&docs).Error
	return docs, err
}
