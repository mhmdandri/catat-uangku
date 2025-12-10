package users

import (
	"errors"

	userprofile "catatan-keuangan/modules/user_profile"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	FindAll() ([]User, error)
	FindByID(ID uuid.UUID) (User, error)
	Create(userRequest UserRequest) (User, error)
	Update(ID uuid.UUID, userRequest UserRequest) (User, error)
	Delete(ID uuid.UUID) (User, error)
	ChangePassword(userID uuid.UUID, req ChangePasswordRequest) error
}

type service struct {
	repository     Repository
	profileService userprofile.Service
}

func NewService(repository Repository, profileService userprofile.Service) *service {
	return &service{repository: repository, profileService: profileService}
}

func (s *service) FindAll() ([]User, error) {
	return s.repository.FindAll()
}

func (s *service) FindByID(ID uuid.UUID) (User, error) {
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
	newUser, err := s.repository.Create(user)
	if err != nil {
		return User{}, err
	}
	if s.profileService != nil {
		profile, err := s.profileService.EnsureDefaultProfile(newUser.ID, newUser.Name, newUser.Email)
		if err != nil {
			return User{}, err
		}
		newUser.Profile = profile
	}
	return newUser, nil
}

func (s *service) Update(ID uuid.UUID, userRequest UserRequest) (User, error) {
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

func (s *service) Delete(ID uuid.UUID) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return User{}, err
	}
	return s.repository.Delete(user)
}

func (s *service) ChangePassword(userID uuid.UUID, req ChangePasswordRequest) error {
	user, err := s.repository.FindByID(userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("password lama salah")
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("konfirmasi password tidak sama")
	}
	if req.NewPassword == req.OldPassword {
		return errors.New("password baru tidak boleh sama dengan password lama")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal membuat password")
	}
	user.Password = string(hashPassword)
	_, err = s.repository.Update(user)
	return err
}
