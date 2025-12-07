package attachments

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(attachment Attachment) (Attachment, error)
	FindByTransactionID(transactionID uuid.UUID) ([]Attachment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(attachment Attachment) (Attachment, error) {
	err := r.db.Create(&attachment).Error
	return attachment, err
}

func (r *repository) FindByTransactionID(transactionID uuid.UUID) ([]Attachment, error) {
	var attach []Attachment
	err := r.db.Where("transaction_id = ?", transactionID).Find(&attach).Error
	return attach, err
}
