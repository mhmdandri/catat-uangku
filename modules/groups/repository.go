package groups

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(group Group) (Group, error)
	FindByID(ID uuid.UUID) (Group, error)
	FindGroupByUserID(userID uuid.UUID) ([]Group, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(group Group) (Group, error) {
	err := r.db.Create(&group).Error
	return group, err
}

func (r *repository) FindByID(ID uuid.UUID) (Group, error) {
	var group Group
	err := r.db.Preload("GroupMembers").First(&group, "id = ?", ID).Error
	return group, err
}

func (r *repository) FindGroupByUserID(userID uuid.UUID) ([]Group, error) {
	var groups []Group
	err := r.db.Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ? AND group_members.is_active = true", userID).
		Find(&groups).Error
	return groups, err
}
