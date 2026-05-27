package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByID(id uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) UpdateName(id uuid.UUID, name string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("name", name).Error
}

func (r *UserRepo) UpdatePassword(id uuid.UUID, hash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password_hash", hash).Error
}

func (r *UserRepo) IncrementUploads(id uuid.UUID) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		UpdateColumn("uploads_this_month", gorm.Expr("uploads_this_month + 1")).Error
}

func (r *UserRepo) ResetUploads() error {
	return r.db.Model(&models.User{}).
		Where("1 = 1").
		Updates(map[string]interface{}{"uploads_this_month": 0, "uploads_reset_at": gorm.Expr("NOW()")}).Error
}
