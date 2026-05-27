package service

import (
	"errors"

	"github.com/google/uuid"

	"quikram-service/internal/models"
	"quikram-service/internal/repository"
	"quikram-service/internal/validator"
)

type SubjectService struct {
	subjectRepo *repository.SubjectRepo
}

func NewSubjectService(subjectRepo *repository.SubjectRepo) *SubjectService {
	return &SubjectService{subjectRepo: subjectRepo}
}

func (s *SubjectService) List(userID uuid.UUID) ([]models.Subject, error) {
	return s.subjectRepo.FindByUserID(userID)
}

func (s *SubjectService) Create(userID uuid.UUID, title, category string) (*models.Subject, error) {
	if err := validator.SubjectTitle(title); err != nil {
		return nil, err
	}

	count, err := s.subjectRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	if count >= 3 {
		return nil, errors.New("free plan limit reached: max 3 subjects")
	}

	subject := &models.Subject{
		UserID:   userID,
		Title:    title,
		Category: category,
	}
	if err := s.subjectRepo.Create(subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *SubjectService) GetByID(id, userID uuid.UUID) (*models.Subject, error) {
	subject, err := s.subjectRepo.FindByID(id)
	if err != nil || subject == nil {
		return nil, errors.New("subject not found")
	}
	if subject.UserID != userID {
		return nil, errors.New("access denied")
	}
	return subject, nil
}

func (s *SubjectService) Update(id, userID uuid.UUID, title, category string) error {
	subject, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	return s.subjectRepo.Update(subject.ID, title, category)
}

func (s *SubjectService) Delete(id, userID uuid.UUID) error {
	subject, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	return s.subjectRepo.Delete(subject.ID)
}
