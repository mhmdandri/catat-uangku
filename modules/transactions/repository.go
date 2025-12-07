package transactions

import "gorm.io/gorm"

type Repository interface {
	Create(transaction Transactions) (Transactions, error)
	FindByID(ID int) (Transactions, error)
	FindAll() ([]Transactions, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(transaction Transactions) (Transactions, error) {
	err := r.db.Create(&transaction).Error
	return transaction, err
}

func (r *repository) FindByID(ID int) (Transactions, error) {
	var transaction Transactions
	err := r.db.Preload("TransactionLines").Preload("Attachments").First(&transaction, ID).Error
	return transaction, err
}

func (r *repository) FindAll() ([]Transactions, error) {
	var transactions []Transactions
	err := r.db.Preload("TransactionLines").Preload("Attachments").Find(&transactions).Error
	return transactions, err
}
