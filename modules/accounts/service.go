package accounts

import "github.com/google/uuid"

type Service interface {
	Create(accountRequest AccountRequest) (Account, error)
	FindByID(ID uuid.UUID) (Account, error)
	Update(ID uuid.UUID, accountRequest AccountUpdateRequest) (Account, error)
	Delete(ID uuid.UUID) (Account, error)
	FindByUserID(userID uuid.UUID) (Account, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) Create(accountRequest AccountRequest) (Account, error) {
	account := Account{
		OwnerUserID: accountRequest.OwnerUserID,
		GroupID:     accountRequest.GroupID,
		Name:        accountRequest.Name,
		Type:        accountRequest.Type,
		Currency:    accountRequest.Currency,
		Scope:       accountRequest.Scope,
		IsShared:    accountRequest.IsShared,
		IsActive:    accountRequest.IsActive,
	}
	newAccount, err := s.repository.Create(account)
	return newAccount, err
}

func (s *service) FindByID(ID uuid.UUID) (Account, error) {
	user, err := s.repository.FindByID(ID)
	return user, err
}

func (s *service) Update(ID uuid.UUID, accountRequest AccountUpdateRequest) (Account, error) {
	account, err := s.repository.FindByID(ID)
	if err != nil {
		return Account{}, err
	}
	account.Name = accountRequest.Name
	account.Type = accountRequest.Type
	account.IsShared = accountRequest.IsShared
	account.IsActive = accountRequest.IsActive
	return s.repository.Update(account)
}

func (s *service) Delete(ID uuid.UUID) (Account, error) {
	account, err := s.repository.FindByID(ID)
	if err != nil {
		return Account{}, err
	}
	return s.repository.Delete(account)
}

func (s *service) FindByUserID(userID uuid.UUID) (Account, error) {
	userAccount, err := s.repository.FindByUserID(userID)
	if err != nil {
		return Account{}, err
	}
	return userAccount, err
}
