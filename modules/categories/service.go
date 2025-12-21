package categories

import (
	"catatan-keuangan/modules/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(userID uuid.UUID, categoryRequest CategoryRequest) (Category, error)
	GetAllCategories(userID uuid.UUID) ([]Category, error)
}

type service struct {
	repository Repository
	db         *gorm.DB
}

func NewService(repository Repository, db *gorm.DB) *service {
	return &service{repository, db}
}

func (s *service) Create(userID uuid.UUID, categoryRequest CategoryRequest) (Category, error) {
	category := Category{
		Name:  categoryRequest.Name,
		Type:  categoryRequest.Type,
		Color: categoryRequest.Color,
		Icon:  categoryRequest.Icon,
	}
	var groupID *uuid.UUID
	if categoryRequest.GroupID != nil {
		if err := s.ensureExists("groups", *categoryRequest.GroupID); err != nil {
			return Category{}, err
		}
		if err := common.EnsureGroupMember(s.db, *categoryRequest.GroupID, userID); err != nil {
			return Category{}, err
		}
		groupID = categoryRequest.GroupID
	}
	category.GroupID = groupID
	if groupID == nil {
		category.OwnerUserID = &userID
	}
	newCategory, err := s.repository.Create(category)
	return newCategory, err
}

func (s *service) GetAllCategories(userID uuid.UUID) ([]Category, error) {
	var categories []Category
	err := s.db.
		Model(&Category{}).
		Joins("LEFT JOIN group_members gm ON gm.group_id = categories.group_id AND gm.user_id = ? AND gm.is_active = true", userID).
		Where("categories.owner_user_id = ? OR (categories.group_id IS NULL AND categories.owner_user_id IS NULL) OR gm.user_id IS NOT NULL", userID).
		Find(&categories).Error
	return categories, err
}

func (s *service) ensureExists(table string, id uuid.UUID) error {
	var count int64
	if err := s.db.Table(table).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
