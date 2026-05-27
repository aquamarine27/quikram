package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type TokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) Create(t *models.RefreshToken) error {
	return r.db.Create(t).Error
}

func (r *TokenRepo) FindByHash(hash string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *TokenRepo) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *TokenRepo) DeleteByHash(hash string) error {
	return r.db.Where("token_hash = ?", hash).Delete(&models.RefreshToken{}).Error
}
