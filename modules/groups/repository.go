package groups

import "gorm.io/gorm"

type Repository interface {
	Create(group Group) (Group, error)
	FindByID(ID int) (Group, error)
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

func (r *repository) FindByID(ID int) (Group, error) {
	var group Group
	err := r.db.Preload("GroupMembers").First(&group, ID).Error
	return group, err
}
