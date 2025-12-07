package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	FindAll() ([]User, error)
	FindByID(ID int) (User, error)
	Create(userRequest UserRequest) (User, error)
	Update(ID int, userRequest UserRequest) (User, error)
	Delete(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAll() ([]User, error) {
	return s.repository.FindAll()
}

func (s *service) FindByID(ID int) (User, error) {
	return s.repository.FindByID(ID)
}

func (s *service) Create(userRequest UserRequest) (User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New("gagal membuat password")
	}
	user := User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: string(hashPassword),
	}
	return s.repository.Create(user)
}

func (s *service) Update(ID int, userRequest UserRequest) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return User{}, err
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New("gagal membuat password")
	}
	user.Name = userRequest.Name
	user.Email = userRequest.Email
	user.Password = string(hashPassword)
	return s.repository.Update(user)
}

func (s *service) Delete(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return User{}, err
	}
	return s.repository.Delete(user)
}
