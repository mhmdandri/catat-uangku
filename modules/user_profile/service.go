package userprofile

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetProfileByUserID(userID uuid.UUID) (UserProfile, error)
	Create(profile UserProfile) (UserProfile, error)
	Update(userID uuid.UUID, req UpdateProfileRequest) (UserProfile, error)
	EnsureDefaultProfile(userID uuid.UUID, name, email string) (UserProfile, error)
	UpdateAvatar(userID uuid.UUID, avatarPath string) (UserProfile, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetProfileByUserID(userID uuid.UUID) (UserProfile, error) {
	profile, err := s.repository.GetProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserProfile{}, fmt.Errorf("profil user tidak ditemukan: %w", err)
		}
		return UserProfile{}, fmt.Errorf("gagal mengambil profil user: %w", err)
	}
	return profile, err
}

func (s *service) Create(profile UserProfile) (UserProfile, error) {
	if strings.TrimSpace(profile.FirstName) == "" {
		profile.FirstName = "Pengguna"
	}
	profile, err := s.repository.Create(profile)
	if err != nil {
		return UserProfile{}, fmt.Errorf("gagal membuat profil user: %w", err)
	}
	return profile, err
}

func (s *service) EnsureDefaultProfile(userID uuid.UUID, name, email string) (UserProfile, error) {
	profile, err := s.repository.GetProfileByUserID(userID)
	if err == nil {
		return profile, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return UserProfile{}, fmt.Errorf("gagal mengambil profil user: %w", err)
	}
	first, last := splitName(name)
	profile = UserProfile{
		UserID:    userID,
		FirstName: first,
		LastName:  last,
		Email:     email,
	}
	return s.Create(profile)
}

func splitName(name string) (string, *string) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "Pengguna", nil
	}
	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return "Pengguna", nil
	}
	if len(parts) == 1 {
		return parts[0], nil
	}
	first := parts[0]
	last := strings.Join(parts[1:], " ")
	return first, &last
}

func (s *service) Update(userID uuid.UUID, req UpdateProfileRequest) (UserProfile, error) {
	profile, err := s.repository.GetProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserProfile{}, fmt.Errorf("profil user tidak ditemukan: %w", err)
		}
		return UserProfile{}, fmt.Errorf("gagal mengambil profil user: %w", err)
	}
	if req.FirstName != nil {
		profile.FirstName = strings.TrimSpace(*req.FirstName)
		if profile.FirstName == "" {
			profile.FirstName = "Pengguna"
		}
	}
	if req.LastName != nil {
		last := strings.TrimSpace(*req.LastName)
		if last == "" {
			profile.LastName = nil
		} else {
			profile.LastName = &last
		}
	}
	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		if phone == "" {
			profile.Phone = nil
		} else {
			profile.Phone = &phone
		}
	}
	if req.Address != nil {
		addr := strings.TrimSpace(*req.Address)
		if addr == "" {
			profile.Address = nil
		} else {
			profile.Address = &addr
		}
	}
	if req.Bio != nil {
		bio := strings.TrimSpace(*req.Bio)
		if bio == "" {
			profile.Bio = nil
		} else {
			profile.Bio = &bio
		}
	}
	if req.Birthdate != nil {
		profile.Birthdate = req.Birthdate
		profile.Age = calculateAge(req.Birthdate)
	}
	if req.Age != nil {
		profile.Age = req.Age
	}
	if req.AvatarURL != nil {
		avatar := strings.TrimSpace(*req.AvatarURL)
		if avatar == "" {
			profile.AvatarURL = nil
		} else {
			profile.AvatarURL = &avatar
		}
	}
	updated, err := s.repository.Update(profile)
	if err != nil {
		return UserProfile{}, fmt.Errorf("gagal memperbarui profil user: %w", err)
	}
	return updated, nil
}

func (s *service) UpdateAvatar(userID uuid.UUID, avatarPath string) (UserProfile, error) {
	profile, err := s.repository.GetProfileByUserID(userID)
	if err != nil {
		return UserProfile{}, fmt.Errorf("profil user tidak ditemukan: %w", err)
	}
	trimmed := strings.TrimSpace(avatarPath)
	if trimmed == "" {
		return UserProfile{}, errors.New("path avatar kosong")
	}
	profile.AvatarURL = &trimmed
	updated, err := s.repository.Update(profile)
	if err != nil {
		return UserProfile{}, fmt.Errorf("gagal memperbarui avatar: %w", err)
	}
	return updated, nil
}

func calculateAge(birthdate *time.Time) *int {
	if birthdate == nil {
		return nil
	}
	now := time.Now()
	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}
	if years < 0 {
		years = 0
	}
	return &years
}
