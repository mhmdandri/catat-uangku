package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshRepository interface {
	Create(refreshToken RefreshToken) (RefreshToken, error)
	Revoke(id uuid.UUID) error
	FindValidByHash(hash string) (RefreshToken, error)
}

type refreshRepository struct {
	db *gorm.DB
}

func NewRefreshRepository(db *gorm.DB) *refreshRepository {
	return &refreshRepository{db}
}

func (r *refreshRepository) Create(refreshToken RefreshToken) (RefreshToken, error) {
	return refreshToken, r.db.Create(&refreshToken).Error
}

func (r *refreshRepository) Revoke(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&RefreshToken{}).Where("id = ?", id).Update("revoked_at", now).Error
}

func (r *refreshRepository) FindValidByHash(hash string) (RefreshToken, error) {
	var refreshToken RefreshToken
	err := r.db.Where("token = ? AND revoked_at IS NULL AND expires_at > ?", hash, time.Now()).First(&refreshToken).Error
	return refreshToken, err
}
