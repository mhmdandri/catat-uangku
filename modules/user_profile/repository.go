package userprofile

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetProfileByUserID(userID uuid.UUID) (UserProfile, error)
	Create(profile UserProfile) (UserProfile, error)
	Update(profile UserProfile) (UserProfile, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetProfileByUserID(userID uuid.UUID) (UserProfile, error) {
	var profile UserProfile
	err := r.db.First(&profile, "user_id = ?", userID).Error
	return profile, err
}

func (r *repository) Create(profile UserProfile) (UserProfile, error) {
	err := r.db.Create(&profile).Error
	return profile, err
}

func (r *repository) Update(profile UserProfile) (UserProfile, error) {
	err := r.db.Save(&profile).Error
	return profile, err
}
