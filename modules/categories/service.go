package categories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(categoryRequest CategoryRequest) (Category, error)
}

type service struct {
	repository Repository
	db         *gorm.DB
}

func NewService(repository Repository, db *gorm.DB) *service {
	return &service{repository, db}
}

func (s *service) Create(categoryRequest CategoryRequest) (Category, error) {
	if categoryRequest.GroupID == nil && categoryRequest.OwnerUserID == nil {
		return Category{}, errors.New("group_id atau owner_user_id wajib diisi")
	}
	var groupID *uuid.UUID
	if categoryRequest.GroupID != nil {
		var count int64
		if err := s.db.Table("groups").Where("id = ?", *categoryRequest.GroupID).Count(&count).Error; err != nil {
			return Category{}, err
		}
		if count == 0 {
			return Category{}, gorm.ErrRecordNotFound
		}
		groupID = categoryRequest.GroupID
	}
	var ownerUserID *uuid.UUID
	if categoryRequest.OwnerUserID != nil {
		if err := s.ensureExists("users", *categoryRequest.OwnerUserID); err != nil {
			return Category{}, err
		}
		ownerUserID = categoryRequest.OwnerUserID
	}
	category := Category{
		GroupID:     groupID,
		OwnerUserID: ownerUserID,
		Name:        categoryRequest.Name,
		Type:        categoryRequest.Type,
	}
	newCategory, err := s.repository.Create(category)
	return newCategory, err
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
