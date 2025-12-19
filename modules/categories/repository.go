package categories

import "gorm.io/gorm"

type Repository interface {
	Create(category Category) (Category, error)
	GetAllCategories() ([]Category, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(category Category) (Category, error) {
	err := r.db.Create(&category).Error
	return category, err
}

func (r *repository) GetAllCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Find(&categories).Error
	return categories, err
}
