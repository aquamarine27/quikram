package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"quikram-service/internal/models"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (r *ChatRepo) Save(userID uuid.UUID, role, content string) (*models.ChatMessage, error) {
	msg := &models.ChatMessage{
		UserID:  userID,
		Role:    role,
		Content: content,
	}
	err := r.db.Create(msg).Error
	return msg, err
}

func (r *ChatRepo) FindByUserID(userID uuid.UUID, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Where("user_id = ?", userID).
		Order("created_at ASC").Limit(limit).Find(&messages).Error
	return messages, err
}

func (r *ChatRepo) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.ChatMessage{}).Error
}

func (r *ChatRepo) FindLastByUserID(userID uuid.UUID, limit int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Find(&messages).Error
	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, err
}
