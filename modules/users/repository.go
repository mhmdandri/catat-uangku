package users

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]User, error)
	FindByID(ID uuid.UUID) (User, error)
	FindByEmail(email string) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Preload("Accounts").Preload("GroupMembers").Preload("Profile").Find(&users).Error
	return users, err
}

func (r *repository) FindByID(ID uuid.UUID) (User, error) {
	var user User
	err := r.db.Preload("Accounts").Preload("GroupMembers").Preload("Profile").First(&user, "id = ?", ID).Error
	return user, err
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Preload("Profile").First(&user, "email = ?", email).Error
	return user, err
}

func (r *repository) Create(user User) (User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *repository) Delete(user User) (User, error) {
	err := r.db.Delete(&user).Error
	return user, err
}
