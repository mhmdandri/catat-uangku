package accounts

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(account Account) (Account, error)
	Update(account Account) (Account, error)
	FindByID(ID uuid.UUID) (Account, error)
	Delete(account Account) (Account, error)
	FindByUserID(userID uuid.UUID) (Account, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(account Account) (Account, error) {
	err := r.db.Create(&account).Error
	return account, err
}

func (r *repository) FindByID(ID uuid.UUID) (Account, error) {
	var account Account
	err := r.db.First(&account, "id = ?", ID).Error
	return account, err
}

func (r *repository) Update(account Account) (Account, error) {
	err := r.db.Save(&account).Error
	return account, err
}

func (r *repository) Delete(account Account) (Account, error) {
	err := r.db.Delete(&account).Error
	return account, err
}

func (r *repository) FindByUserID(userID uuid.UUID) (Account, error) {
	var account Account
	err := r.db.Where("owner_user_id = ?", userID).First(&account).Error
	return account, err
}
