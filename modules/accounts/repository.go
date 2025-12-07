package accounts

import "gorm.io/gorm"

type Repository interface {
	Create(account Account) (Account, error)
	Update(account Account) (Account, error)
	FindByID(ID int) (Account, error)
	Delete(account Account) (Account, error)
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

func (r *repository) FindByID(ID int) (Account, error) {
	var account Account
	err := r.db.First(&account, ID).Error
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
