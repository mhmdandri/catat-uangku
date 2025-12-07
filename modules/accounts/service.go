package accounts

type Service interface {
	Create(accountRequest AccountRequest) (Account, error)
	FindByID(ID int) (Account, error)
	Update(ID int, accountRequest AccountRequest) (Account, error)
	Delete(ID int) (Account, error)
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

func (s *service) FindByID(ID int) (Account, error) {
	user, err := s.repository.FindByID(ID)
	return user, err
}

func (s *service) Update(ID int, accountRequest AccountRequest) (Account, error) {
	account, err := s.repository.FindByID(ID)
	if err != nil {
		return Account{}, err
	}
	account.OwnerUserID = accountRequest.OwnerUserID
	account.GroupID = accountRequest.GroupID
	account.Name = accountRequest.Name
	account.Type = accountRequest.Type
	account.Currency = accountRequest.Currency
	account.Scope = accountRequest.Scope
	account.IsShared = accountRequest.IsShared
	account.IsActive = accountRequest.IsActive
	return s.repository.Update(account)
}

func (s *service) Delete(ID int) (Account, error) {
	account, err := s.repository.FindByID(ID)
	if err != nil {
		return Account{}, err
	}
	return s.repository.Delete(account)
}
