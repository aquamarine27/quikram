package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type QuizRepo struct {
	db *gorm.DB
}

func NewQuizRepo(db *gorm.DB) *QuizRepo {
	return &QuizRepo{db: db}
}

func (r *QuizRepo) Create(q *models.Quiz) error {
	return r.db.Create(q).Error
}

func (r *QuizRepo) FindBySubjectID(subjectID uuid.UUID) ([]models.Quiz, error) {
	var quizzes []models.Quiz
	err := r.db.Where("subject_id = ?", subjectID).Order("created_at DESC").Find(&quizzes).Error
	return quizzes, err
}

func (r *QuizRepo) FindByID(id uuid.UUID) (*models.Quiz, error) {
	var q models.Quiz
	err := r.db.First(&q, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *QuizRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Quiz{}, "id = ?", id).Error
}

func (r *QuizRepo) CreateQuestions(questions []models.Question) error {
	return r.db.Create(&questions).Error
}

func (r *QuizRepo) FindQuestionsByQuizID(quizID uuid.UUID) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Where("quiz_id = ?", quizID).Order("order_index").Find(&questions).Error
	return questions, err
}

func (r *QuizRepo) DeleteQuestionsByQuizID(quizID uuid.UUID) error {
	return r.db.Where("quiz_id = ?", quizID).Delete(&models.Question{}).Error
}
