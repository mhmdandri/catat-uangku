package transactions

import (
	"catatan-keuangan/modules/attachments"
	transactionlines "catatan-keuangan/modules/transaction_lines"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTransactionsNotFound = errors.New("transactions not found")

type Service interface {
	Create(transactionReq TransactionRequest) (Transactions, error)
	FindByID(ID uuid.UUID) (Transactions, error)
	FindAll() ([]Transactions, error)
	GetTransactionByAccountID(accountID uuid.UUID) ([]Transactions, error)
	GetTransactionByUserID(userID uuid.UUID) ([]Transactions, error)
}
type service struct {
	repository        Repository
	db                *gorm.DB
	attachmentService attachments.Service
}

func NewService(repository Repository, db *gorm.DB, attachmentSvc attachments.Service) *service {
	return &service{
		repository:        repository,
		db:                db,
		attachmentService: attachmentSvc,
	}
}

func (s *service) Create(TransactionRequest TransactionRequest) (Transactions, error) {
	var newTransaction Transactions
	var timeNow = time.Now()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		groupID, err := s.resolveGroupID(tx, TransactionRequest.Scope, TransactionRequest.GroupID)
		if err != nil {
			return err
		}
		if err := s.ensureExists(tx, "categories", TransactionRequest.CategoryID); err != nil {
			return err
		}
		if err := s.ensureExists(tx, "accounts", TransactionRequest.AccountID); err != nil {
			return err
		}
		if TransactionRequest.TotalAmount <= 0 {
			return errors.New("total_amount harus lebih besar dari 0")
		}

		transaction := Transactions{
			GroupID:         groupID,
			CategoryID:      TransactionRequest.CategoryID,
			CreatedByUserID: TransactionRequest.CreatedByUserID,
			Title:           TransactionRequest.Title,
			Date:            timeNow,
			Type:            TransactionRequest.Type,
			TotalAmount:     TransactionRequest.TotalAmount,
			Scope:           TransactionRequest.Scope,
			Description:     TransactionRequest.Description,
		}
		txRepo := NewRepository(tx)
		createdTransaction, err := txRepo.Create(transaction)
		if err != nil {
			return err
		}
		SavedAtt, err := s.attachmentService.SaveFiles(tx, createdTransaction.ID, TransactionRequest.Attachments)
		if err != nil {
			return err
		}
		createdTransaction.Attachments = SavedAtt
		var debit float64 = 0
		var credit float64 = 0
		switch transaction.Type {
		case "income":
			credit = transaction.TotalAmount
		case "expense":
			debit = transaction.TotalAmount
		}
		createTransactionLine := transactionlines.TransactionLine{
			TransactionID: createdTransaction.ID,
			AccountID:     TransactionRequest.AccountID,
			Debit:         debit,
			Credit:        credit,
			Note:          transaction.Description,
		}
		if err := tx.Create(&createTransactionLine).Error; err != nil {
			return err
		}
		
		// Reload transaction dengan relasi lengkap
		if err := tx.
			Preload("Category").
			Preload("TransactionLines").
			Preload("TransactionLines.Accounts").
			Preload("Attachments").
			First(&createdTransaction, "id = ?", createdTransaction.ID).Error; err != nil {
			return err
		}
		
		newTransaction = createdTransaction
		return nil
	})
	return newTransaction, err
}

func (s *service) FindByID(ID uuid.UUID) (Transactions, error) {
	transaction, err := s.repository.FindByID(ID)
	return transaction, err
}

func (s *service) FindAll() ([]Transactions, error) {
	transactions, err := s.repository.FindAll()
	return transactions, err
}

func (s *service) GetTransactionByAccountID(accountID uuid.UUID) ([]Transactions, error) {
	transactions, err := s.repository.GetTransactionByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *service) GetTransactionByUserID(userID uuid.UUID) ([]Transactions, error) {
	transactions, err := s.repository.GetTransactionByUserID(userID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *service) resolveGroupID(tx *gorm.DB, scope string, groupID *uuid.UUID) (*uuid.UUID, error) {
	switch scope {
	case "group":
		if groupID == nil {
			return nil, errors.New("group_id wajib diisi untuk scope group")
		}
		if err := s.ensureExists(tx, "groups", *groupID); err != nil {
			return nil, err
		}
		return groupID, nil
	case "personal":
		return nil, nil
	default:
		return nil, errors.New("scope tidak valid")
	}
}

func (s *service) ensureExists(tx *gorm.DB, table string, id uuid.UUID) error {
	var count int64
	if err := tx.Table(table).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
