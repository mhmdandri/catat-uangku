package transactions

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(transaction Transactions) (Transactions, error)
	FindByID(ID uuid.UUID) (Transactions, error)
	FindAllByUserID(userID uuid.UUID, startDate, endDate time.Time) ([]Transactions, error)
	Delete(ID uuid.UUID) error
	Update(ID uuid.UUID, transaction Transactions) (Transactions, error)
	GetTransactionByAccountID(accountID uuid.UUID, startDate, endDate time.Time) ([]Transactions, error)
	GetTransactionByUserID(userID uuid.UUID, startDate, endDate time.Time) ([]Transactions, error)
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
	err := r.db.
		Preload("Category").
		Preload("TransactionLines").
		Preload("TransactionLines.Accounts").
		Preload("Attachments").
		First(&transaction, "id = ?", ID).Error
	return transaction, err
}

func (r *repository) FindAllByUserID(userID uuid.UUID, startDate, endDate time.Time) ([]Transactions, error) {
	var transactions []Transactions
	err := r.db.
		Preload("Category").
		Preload("TransactionLines").
		Preload("TransactionLines.Accounts").
		Preload("Attachments").
		Joins("LEFT JOIN group_members gm ON gm.group_id = transactions.group_id AND gm.user_id = ? AND gm.is_active = true", userID).
		Where("transactions.date >= ? AND transactions.date < ?", startDate, endDate).
		Where("(transactions.scope = ? AND transactions.created_by_user_id = ?) OR (transactions.scope = ? AND gm.user_id IS NOT NULL)", "personal", userID, "group").
		Order("transactions.created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *repository) Update(ID uuid.UUID, transaction Transactions) (Transactions, error) {
	err := r.db.Model(&Transactions{}).
		Where("id = ?", ID).
		Select("group_id", "category_id", "title", "type", "total_amount", "scope", "description").
		Updates(&transaction).Error
	return transaction, err
}

func (r *repository) GetTransactionByAccountID(accountID uuid.UUID, startDate, endDate time.Time) ([]Transactions, error) {
	var transactions []Transactions
	err := r.db.
		Preload("Category").
		Preload("TransactionLines").
		Preload("TransactionLines.Accounts").
		Preload("Attachments").
		Joins("JOIN transaction_lines tl ON tl.transaction_id = transactions.id").
		Where("tl.account_id = ? AND transactions.date >= ? AND transactions.date < ?", accountID, startDate, endDate).
		Select("transactions.*").
		Order("created_at desc").
		Find(&transactions).Error
	return transactions, err
}

func (r *repository) GetTransactionByUserID(userID uuid.UUID, startDate, endDate time.Time) ([]Transactions, error) {
	var transactions []Transactions
	err := r.db.
		Preload("Category").
		Preload("TransactionLines").
		Preload("TransactionLines.Accounts").
		Preload("Attachments").
		Where("created_by_user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *repository) Delete(ID uuid.UUID) error {
	return r.db.Delete(&Transactions{}, "id = ?", ID).Error
}
