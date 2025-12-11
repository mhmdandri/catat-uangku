package transactions

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(transaction Transactions) (Transactions, error)
	FindByID(ID uuid.UUID) (Transactions, error)
	FindAll() ([]Transactions, error)
	GetTransactionByAccountID(accountID uuid.UUID) ([]Transactions, error)
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

func (r *repository) FindByID(ID uuid.UUID) (Transactions, error) {
	var transaction Transactions
	err := r.db.Preload("TransactionLines").Preload("Attachments").First(&transaction, "id = ?", ID).Error
	return transaction, err
}

func (r *repository) FindAll() ([]Transactions, error) {
	var transactions []Transactions
	err := r.db.Preload("TransactionLines").Preload("Attachments").Find(&transactions).Error
	return transactions, err
}

func (r *repository) GetTransactionByAccountID(accountID uuid.UUID) ([]Transactions, error) {
	var transactions []Transactions
	err := r.db.
		Preload("TransactionLines").
		Preload("Attachments").
		Joins("JOIN transaction_lines tl ON tl.transaction_id = transactions.id").
		Where("tl.account_id = ?", accountID).
		Select("transactions.*").
		Find(&transactions).Error
	return transactions, err
}
