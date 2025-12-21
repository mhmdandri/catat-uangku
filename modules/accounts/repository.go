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
	FindByUserID(userID uuid.UUID) ([]Account, error)
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
	err := r.db.Model(&Account{}).
		Scopes(WithBalance).
		Where("accounts.id = ?", ID).
		First(&account).Error
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

func (r *repository) FindByUserID(userID uuid.UUID) ([]Account, error) {
	var accounts []Account
	err := r.db.Model(&Account{}).
		Scopes(WithBalance).
		Joins("LEFT JOIN group_members gm ON gm.group_id = accounts.group_id AND gm.user_id = ? AND gm.is_active = true", userID).
		Where("accounts.owner_user_id = ?", userID).
		Or("gm.user_id IS NOT NULL AND accounts.scope = ? AND accounts.is_shared = ?", "group", true).
		Order("accounts.created_at ASC").
		Find(&accounts).Error
	return accounts, err
}
