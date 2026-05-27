package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type AttemptRepo struct {
	db *gorm.DB
}

func NewAttemptRepo(db *gorm.DB) *AttemptRepo {
	return &AttemptRepo{db: db}
}

func (r *AttemptRepo) Create(a *models.QuizAttempt) error {
	return r.db.Create(a).Error
}

func (r *AttemptRepo) Update(a *models.QuizAttempt) error {
	return r.db.Save(a).Error
}

func (r *AttemptRepo) FindByQuizID(quizID uuid.UUID) ([]models.QuizAttempt, error) {
	var attempts []models.QuizAttempt
	err := r.db.Where("quiz_id = ?", quizID).Order("started_at DESC").Find(&attempts).Error
	return attempts, err
}

func (r *AttemptRepo) FindByID(id uuid.UUID) (*models.QuizAttempt, error) {
	var a models.QuizAttempt
	err := r.db.First(&a, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AttemptRepo) GetUserStats(userID uuid.UUID) (totalTests int64, avgScore float64, err error) {
	err = r.db.Model(&models.QuizAttempt{}).
		Where("user_id = ? AND completed_at IS NOT NULL", userID).
		Select("COUNT(*), COALESCE(AVG(score), 0)").
		Row().Scan(&totalTests, &avgScore)
	return
}
