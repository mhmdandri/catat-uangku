package transactions

import (
	"catatan-keuangan/modules/accounts"
	"catatan-keuangan/modules/attachments"
	"catatan-keuangan/modules/categories"
	"catatan-keuangan/modules/common"
	transactionlines "catatan-keuangan/modules/transaction_lines"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTransactionsNotFound = errors.New("transactions not found")

type Service interface {
	Create(userID uuid.UUID, transactionReq TransactionRequest) (Transactions, error)
	FindByID(userID, ID uuid.UUID) (Transactions, error)
	FindAll(userID uuid.UUID) ([]Transactions, error)
	Delete(userID, ID uuid.UUID) error
	Update(userID, ID uuid.UUID, transactionReq TransactionRequest) (Transactions, error)
	GetTransactionByAccountID(userID, accountID uuid.UUID) ([]Transactions, error)
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

func (s *service) Create(userID uuid.UUID, transactionRequest TransactionRequest) (Transactions, error) {
	var newTransaction Transactions
	var timeNow = time.Now()
	err := s.db.Transaction(func(tx *gorm.DB) error {
		groupID, err := s.resolveGroupID(tx, transactionRequest.Scope, transactionRequest.GroupID)
		if err != nil {
			return err
		}
		if groupID != nil {
			if err := common.EnsureGroupMember(tx, *groupID, userID); err != nil {
				return err
			}
		}
		if err := s.ensureCategoryAccess(tx, transactionRequest.CategoryID, userID, groupID); err != nil {
			return err
		}
		account, err := s.getAccountForUser(tx, transactionRequest.AccountID, userID)
		if err != nil {
			return err
		}
		if err := s.ensureAccountScope(account, transactionRequest.Scope, groupID); err != nil {
			return err
		}
		if transactionRequest.TotalAmount <= 0 {
			return errors.New("total_amount harus lebih besar dari 0")
		}

		transaction := Transactions{
			GroupID:         groupID,
			CategoryID:      transactionRequest.CategoryID,
			CreatedByUserID: userID,
			Title:           transactionRequest.Title,
			Date:            timeNow,
			Type:            transactionRequest.Type,
			TotalAmount:     transactionRequest.TotalAmount,
			Scope:           transactionRequest.Scope,
			Description:     transactionRequest.Description,
		}
		txRepo := NewRepository(tx)
		createdTransaction, err := txRepo.Create(transaction)
		if err != nil {
			return err
		}
		SavedAtt, err := s.attachmentService.SaveFiles(tx, createdTransaction.ID, transactionRequest.Attachments)
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
			AccountID:     transactionRequest.AccountID,
			Debit:         debit,
			Credit:        credit,
			Note:          transaction.Description,
		}
		if err := tx.Create(&createTransactionLine).Error; err != nil {
			return err
		}
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

func (s *service) Delete(userID, ID uuid.UUID) error {
	transaction, err := s.repository.FindByID(ID)
	if err != nil {
		return err
	}
	if err := s.ensureTransactionAccess(s.db, transaction, userID); err != nil {
		return err
	}
	return s.repository.Delete(ID)
}

func (s *service) Update(userID, ID uuid.UUID, transactionReq TransactionRequest) (Transactions, error) {
	var updated Transactions
	err := s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)
		existingTransaction, err := txRepo.FindByID(ID)
		if err != nil {
			return err
		}
		if err := s.ensureTransactionAccess(tx, existingTransaction, userID); err != nil {
			return err
		}
		groupID, err := s.resolveGroupID(tx, transactionReq.Scope, transactionReq.GroupID)
		if err != nil {
			return err
		}
		if groupID != nil {
			if err := common.EnsureGroupMember(tx, *groupID, userID); err != nil {
				return err
			}
		}
		if err := s.ensureCategoryAccess(tx, transactionReq.CategoryID, userID, groupID); err != nil {
			return err
		}
		account, err := s.getAccountForUser(tx, transactionReq.AccountID, userID)
		if err != nil {
			return err
		}
		if err := s.ensureAccountScope(account, transactionReq.Scope, groupID); err != nil {
			return err
		}
		if transactionReq.TotalAmount <= 0 {
			return errors.New("total_amount harus lebih besar dari 0")
		}

		updatedTransaction := Transactions{
			ID:          existingTransaction.ID,
			GroupID:     groupID,
			CategoryID:  transactionReq.CategoryID,
			Title:       transactionReq.Title,
			Type:        transactionReq.Type,
			TotalAmount: transactionReq.TotalAmount,
			Scope:       transactionReq.Scope,
			Description: transactionReq.Description,
		}
		if _, err := txRepo.Update(ID, updatedTransaction); err != nil {
			return err
		}

		var debit float64 = 0
		var credit float64 = 0
		switch transactionReq.Type {
		case "income":
			credit = transactionReq.TotalAmount
		case "expense":
			debit = transactionReq.TotalAmount
		}

		var lines []transactionlines.TransactionLine
		if err := tx.Where("transaction_id = ?", ID).Find(&lines).Error; err != nil {
			return err
		}
		switch len(lines) {
		case 0:
			line := transactionlines.TransactionLine{
				TransactionID: ID,
				AccountID:     transactionReq.AccountID,
				Debit:         debit,
				Credit:        credit,
				Note:          transactionReq.Description,
			}
			if err := tx.Create(&line).Error; err != nil {
				return err
			}
		case 1:
			updates := map[string]interface{}{
				"account_id": transactionReq.AccountID,
				"debit":      debit,
				"credit":     credit,
				"note":       transactionReq.Description,
			}
			if err := tx.Model(&transactionlines.TransactionLine{}).
				Where("id = ?", lines[0].ID).
				Updates(updates).Error; err != nil {
				return err
			}
		default:
			return errors.New("transaksi memiliki lebih dari satu line, update tidak didukung")
		}

		if len(transactionReq.Attachments) > 0 {
			if _, err := s.attachmentService.SaveFiles(tx, ID, transactionReq.Attachments); err != nil {
				return err
			}
		}

		updated, err = txRepo.FindByID(ID)
		return err
	})
	if err != nil {
		return Transactions{}, err
	}
	return updated, nil
}

func (s *service) FindByID(userID, ID uuid.UUID) (Transactions, error) {
	transaction, err := s.repository.FindByID(ID)
	if err != nil {
		return Transactions{}, err
	}
	if err := s.ensureTransactionAccess(s.db, transaction, userID); err != nil {
		return Transactions{}, err
	}
	return transaction, err
}

func (s *service) FindAll(userID uuid.UUID) ([]Transactions, error) {
	transactions, err := s.repository.FindAllByUserID(userID)
	return transactions, err
}

func (s *service) GetTransactionByAccountID(userID, accountID uuid.UUID) ([]Transactions, error) {
	if _, err := s.getAccountForUser(s.db, accountID, userID); err != nil {
		return nil, err
	}
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

func (s *service) ensureCategoryAccess(tx *gorm.DB, categoryID, userID uuid.UUID, groupID *uuid.UUID) error {
	var category categories.Category
	if err := tx.First(&category, "id = ?", categoryID).Error; err != nil {
		return err
	}
	if category.GroupID != nil {
		if groupID == nil || *category.GroupID != *groupID {
			return common.ErrForbidden
		}
		if err := common.EnsureGroupMember(tx, *category.GroupID, userID); err != nil {
			return err
		}
		return nil
	}
	if category.OwnerUserID != nil && *category.OwnerUserID != userID {
		return common.ErrForbidden
	}
	return nil
}

func (s *service) getAccountForUser(tx *gorm.DB, accountID, userID uuid.UUID) (accounts.Account, error) {
	var account accounts.Account
	if err := tx.First(&account, "id = ?", accountID).Error; err != nil {
		return account, err
	}
	if account.Scope == "group" {
		if account.GroupID == nil {
			return account, common.ErrForbidden
		}
		if err := common.EnsureGroupMember(tx, *account.GroupID, userID); err != nil {
			return account, err
		}
		if account.OwnerUserID != userID && !account.IsShared {
			return account, common.ErrForbidden
		}
		return account, nil
	}
	if account.OwnerUserID != userID {
		return account, common.ErrForbidden
	}
	return account, nil
}

func (s *service) ensureAccountScope(account accounts.Account, scope string, groupID *uuid.UUID) error {
	if scope == "personal" {
		if account.Scope != "personal" {
			return common.ErrForbidden
		}
		return nil
	}
	if scope == "group" {
		if account.Scope != "group" {
			return common.ErrForbidden
		}
		if account.GroupID == nil || groupID == nil || *account.GroupID != *groupID {
			return common.ErrForbidden
		}
		return nil
	}
	return common.ErrForbidden
}

func (s *service) ensureTransactionAccess(tx *gorm.DB, transaction Transactions, userID uuid.UUID) error {
	switch transaction.Scope {
	case "personal":
		if transaction.CreatedByUserID != userID {
			return common.ErrForbidden
		}
		return nil
	case "group":
		if transaction.GroupID == nil {
			return common.ErrForbidden
		}
		return common.EnsureGroupMember(tx, *transaction.GroupID, userID)
	default:
		return common.ErrForbidden
	}
}
