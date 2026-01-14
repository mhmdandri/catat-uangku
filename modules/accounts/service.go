package accounts

import (
	"errors"

	"catatan-keuangan/modules/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(userID uuid.UUID, accountRequest AccountRequest) (Account, error)
	FindByID(userID, ID uuid.UUID) (Account, error)
	Update(userID, ID uuid.UUID, accountRequest AccountUpdateRequest) (Account, error)
	Delete(userID, ID uuid.UUID) (Account, error)
	FindByUserID(userID uuid.UUID) ([]Account, error)
}

type service struct {
	repository Repository
	db         *gorm.DB
}

func NewService(repository Repository, db *gorm.DB) *service {
	return &service{repository: repository, db: db}
}

var ErrTransactionExists = errors.New("akun sudah ada transaksi tidak bisa di hapus")

func (s *service) Create(userID uuid.UUID, accountRequest AccountRequest) (Account, error) {
	accountRequest.OwnerUserID = userID
	switch accountRequest.Scope {
	case "personal":
		accountRequest.GroupID = nil
	case "group":
		if accountRequest.GroupID == nil {
			return Account{}, errors.New("group_id wajib diisi untuk scope group")
		}
		if err := common.EnsureGroupMember(s.db, *accountRequest.GroupID, userID); err != nil {
			return Account{}, err
		}
	default:
		return Account{}, errors.New("scope tidak valid")
	}
	firstBalance := 0.0
	if accountRequest.FirstBalance != nil {
		firstBalance = *accountRequest.FirstBalance
	}

	account := Account{
		OwnerUserID:  accountRequest.OwnerUserID,
		GroupID:      accountRequest.GroupID,
		Name:         accountRequest.Name,
		Type:         accountRequest.Type,
		Number:       accountRequest.Number,
		FirstBalance: firstBalance,
		Currency:     accountRequest.Currency,
		Scope:        accountRequest.Scope,
		IsShared:     accountRequest.IsShared,
		IsActive:     accountRequest.IsActive,
	}
	newAccount, err := s.repository.Create(account)
	return newAccount, err
}

func (s *service) FindByID(userID, ID uuid.UUID) (Account, error) {
	acc, err := s.repository.FindByID(ID)
	if err != nil {
		return Account{}, err
	}
	if err := s.ensureAccountAccess(acc, userID); err != nil {
		return Account{}, err
	}
	return acc, err
}

func (s *service) Update(userID, ID uuid.UUID, accountRequest AccountUpdateRequest) (Account, error) {
	account, err := s.repository.FindByID(ID)
	if err != nil {
		return Account{}, err
	}
	if err := s.ensureAccountAccess(account, userID); err != nil {
		return Account{}, err
	}
	account.Name = accountRequest.Name
	account.Type = accountRequest.Type
	account.Number = accountRequest.Number
	account.Currency = accountRequest.Currency
	account.IsShared = accountRequest.IsShared
	account.IsActive = accountRequest.IsActive
	return s.repository.Update(account)
}

func (s *service) Delete(userID, ID uuid.UUID) (Account, error) {
	account, err := s.repository.FindByID(ID)
	if err != nil {
		return Account{}, err
	}
	if err := s.ensureAccountAccess(account, userID); err != nil {
		return Account{}, err
	}
	if err := s.ensureNoTransactions(ID); err != nil {
		return Account{}, err
	}
	return s.repository.Delete(account)
}

func (s *service) FindByUserID(userID uuid.UUID) ([]Account, error) {
	userAccount, err := s.repository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return userAccount, err
}

func (s *service) ensureAccountAccess(account Account, userID uuid.UUID) error {
	if account.Scope == "group" {
		if account.GroupID == nil {
			return common.ErrForbidden
		}
		if err := common.EnsureGroupMember(s.db, *account.GroupID, userID); err != nil {
			return err
		}
		if account.OwnerUserID != userID && !account.IsShared {
			return common.ErrForbidden
		}
		return nil
	}
	if account.OwnerUserID != userID {
		return common.ErrForbidden
	}
	return nil
}

func (s *service) ensureNoTransactions(accountID uuid.UUID) error {
	var count int64
	if err := s.db.Table("transaction_lines").
		Where("account_id = ?", accountID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrTransactionExists
	}
	return nil
}
