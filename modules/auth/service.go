package auth

import (
	"catatan-keuangan/modules/users"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Login(loginRequest LoginRequest) (string, string, users.User, error)
	Refresh(raw string) (string, string, users.User, error)
}

type service struct {
	userRepository    users.Repository
	refreshRepository RefreshRepository
}

func NewService(userRepository users.Repository, refreshRepository RefreshRepository) *service {
	return &service{userRepository, refreshRepository}
}

func generateRefresh() (string, string, error) {
	base := make([]byte, 32)
	if _, err := rand.Read(base); err != nil {
		return "", "", err
	}
	raw := base64.StdEncoding.EncodeToString(base)
	sum := sha256.Sum256([]byte(raw))
	return raw, base64.StdEncoding.EncodeToString(sum[:]), nil
}

func (s *service) Login(loginRequest LoginRequest) (string, string, users.User, error) {
	user, err := s.userRepository.FindByEmail(loginRequest.Email)
	if err != nil {
		return "", "", users.User{}, errors.New("email atau password salah")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil {
		return "", "", users.User{}, errors.New("email atau password salah")
	}
	access, _ := GenerateToken(TokenData{UserID: user.ID, Role: ""})
	rawRefresh, hashRefresh, _ := generateRefresh()
	expires := time.Now().Add(7 * 24 * time.Hour)
	s.refreshRepository.Create(RefreshToken{UserID: user.ID, Token: hashRefresh, ExpiresAt: expires})
	return access, rawRefresh, user, nil
}

func hashRaw(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func (s *service) Refresh(raw string) (string, string, users.User, error) {
	hash := hashRaw(raw) // gunakan sha256 seperti di generateRawRefresh
	rt, err := s.refreshRepository.FindValidByHash(hash)
	if err != nil {
		return "", "", users.User{}, errors.New("refresh token invalid")
	}

	user, err := s.userRepository.FindByID(rt.UserID)
	if err != nil {
		return "", "", users.User{}, err
	}

	if err := s.refreshRepository.Revoke(rt.ID); err != nil {
		return "", "", users.User{}, err
	}

	access, _ := GenerateToken(TokenData{UserID: user.ID, Role: ""})
	newRaw, newHash, _ := generateRefresh()
	s.refreshRepository.Create(RefreshToken{UserID: user.ID, Token: newHash, ExpiresAt: time.Now().Add(7 * 24 * time.Hour), RotatedFrom: &rt.ID})
	return access, newRaw, user, nil
}
