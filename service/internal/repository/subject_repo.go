package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type SubjectRepo struct {
	db *gorm.DB
}

func NewSubjectRepo(db *gorm.DB) *SubjectRepo {
	return &SubjectRepo{db: db}
}

func (r *SubjectRepo) Create(s *models.Subject) error {
	return r.db.Create(s).Error
}

func (r *SubjectRepo) FindByUserID(userID uuid.UUID) ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&subjects).Error
	return subjects, err
}

func (r *SubjectRepo) FindByID(id uuid.UUID) (*models.Subject, error) {
	var s models.Subject
	err := r.db.First(&s, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SubjectRepo) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Subject{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *SubjectRepo) Update(id uuid.UUID, title, category string) error {
	return r.db.Model(&models.Subject{}).Where("id = ?", id).
		Updates(map[string]interface{}{"title": title, "category": category}).Error
}

func (r *SubjectRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Subject{}, "id = ?", id).Error
}
