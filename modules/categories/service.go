package categories

import (
	"errors"

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
	if categoryRequest.GroupID <= 0 && categoryRequest.OwnerUserID <= 0 {
		return Category{}, errors.New("group_id atau owner_user_id wajib diisi")
	}
	var groupID *int
	if categoryRequest.GroupID > 0 {
		var count int64
		if err := s.db.Table("groups").Where("id = ?", categoryRequest.GroupID).Count(&count).Error; err != nil {
			return Category{}, err
		}
		if count == 0 {
			return Category{}, gorm.ErrRecordNotFound
		}
		groupID = &categoryRequest.GroupID
	}
	var ownerUserID *int
	if categoryRequest.OwnerUserID > 0 {
		if err := s.ensureExists("users", categoryRequest.OwnerUserID); err != nil {
			return Category{}, err
		}
		ownerUserID = &categoryRequest.OwnerUserID
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

func (s *service) ensureExists(table string, id int) error {
	var count int64
	if err := s.db.Table(table).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
